package types

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/golang/protobuf/proto"
	hederaproto "github.com/hashgraph/hedera-sdk-go/v2/proto"
	"github.com/hashgraph/hedera-state-proof-verifier-go/internal/constants"
	"github.com/hashgraph/hedera-state-proof-verifier-go/internal/errors"
	"github.com/hashgraph/hedera-state-proof-verifier-go/internal/reader"
	"io"
	"reflect"
)

type record struct {
	Stream
	BodyLength  uint32
	Record      []byte
	Transaction []byte
}

type RecordFile struct {
	Hash            string
	MetadataHash    string
	TransactionsMap map[string]*hederaproto.TransactionID
}

func newRecordFromString(record string) (*record, error) {
	decoded, err := base64.StdEncoding.DecodeString(record)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewReader(decoded)

	return newRecord(buffer)
}

func newRecord(buffer *bytes.Reader) (*record, error) {
	stream, err := NewStream(buffer)
	if err != nil {
		return nil, err
	}

	record := &record{
		Stream: *stream,
	}

	bodyLength, err := record.readBody(buffer)
	if err != nil {
		return nil, err
	}

	record.BodyLength = *bodyLength + stream.BodyLength

	return record, nil
}

func (r *record) readBody(buffer *bytes.Reader) (*uint32, error) {
	recordLength, recordBytes, err := reader.LengthAndBytes(buffer, constants.ByteSize, constants.MaxRecordLength, false)
	if err != nil {
		return nil, err
	}

	transactionLength, transactionBytes, err := reader.LengthAndBytes(buffer, constants.ByteSize, constants.MaxTransactionLength, false)
	if err != nil {
		return nil, err
	}

	r.Record = recordBytes
	r.Transaction = transactionBytes

	totalLength := *recordLength + *transactionLength

	return &totalLength, nil
}

func NewPreV5RecordFile(buffer *bytes.Reader) (*RecordFile, error) {
	bytesToRead := make([]byte, constants.PreV5HeaderLength)
	err := binary.Read(buffer, binary.BigEndian, bytesToRead)
	if err != nil {
		return nil, err
	}

	recordFile := &RecordFile{
		TransactionsMap: make(map[string]*hederaproto.TransactionID),
	}

	for buffer.Len() != 0 {
		marker, err := buffer.ReadByte()
		if err != nil {
			return nil, err
		}
		if marker != constants.RecordFileMarker {
			return nil, errors.ErrorUnsupportedRecordFileMarker
		}

		_, _, err = reader.LengthAndBytes(buffer, constants.ByteSize, constants.MaxTransactionLength, false)
		if err != nil {
			return nil, err
		}
		_, recordBytes, err := reader.LengthAndBytes(buffer, constants.ByteSize, constants.MaxRecordLength, false)
		if err != nil {
			return nil, err
		}

		err = mapSuccessfulTransactions(recordFile.TransactionsMap, recordBytes)
		if err != nil {
			return nil, err
		}
	}

	return recordFile, nil
}

// skip the bytes before the start hash object to read a list of stream objects organized as follows:
//
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |  Start Object Running Hash  |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |    Record Stream Object     |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |    ...                      |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |    Record Stream Object     |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |   End Object Running Hash   |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//
// Note the start object running hash and the end object running hash are of the same type HashObject and
// they have the same classId.
func NewV5RecordFile(buffer *bytes.Reader) (*RecordFile, error) {
	var metadata []byte
	hashOffset := make([]byte, constants.V5StartHashOffset)

	err := binary.Read(buffer, binary.BigEndian, hashOffset)
	if err != nil {
		return nil, err
	}
	metadata = append(metadata, hashOffset...)

	hash, err := NewHash(buffer)
	if err != nil {
		return nil, err
	}

	_, err = buffer.Seek(-int64(hash.BodyLength), io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	hashBytes := make([]byte, hash.BodyLength)
	err = binary.Read(buffer, binary.BigEndian, hashBytes)
	if err != nil {
		return nil, err
	}
	metadata = append(metadata, hashBytes...)

	recordFile := &RecordFile{
		TransactionsMap: make(map[string]*hederaproto.TransactionID),
	}

	// record stream objects are between the start hash object and the end hash object
	var classId int64
	for {
		err := binary.Read(buffer, binary.BigEndian, &classId)
		if err != nil {
			return nil, err
		}

		_, err = buffer.Seek(-constants.LongSize, io.SeekCurrent)
		if err != nil {
			return nil, err
		}

		if classId == hash.ClassId {
			break
		}

		record, err := newRecord(buffer)
		if err != nil {
			return nil, err
		}

		err = mapSuccessfulTransactions(recordFile.TransactionsMap, record.Record)
		if err != nil {
			return nil, err
		}
	}

	endHash, err := NewHash(buffer)
	if err != nil {
		return nil, err
	}

	if buffer.Len() != 0 {
		return nil, errors.ErrorExtraDataInRecordFile
	}

	_, err = buffer.Seek(-int64(endHash.BodyLength), io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	err = binary.Read(buffer, binary.BigEndian, hashBytes)
	if err != nil {
		return nil, err
	}

	metadata = append(metadata, hashBytes...)

	metadataHash := sha512.Sum384(metadata)

	recordFile.MetadataHash = hex.EncodeToString(metadataHash[:])

	return recordFile, nil
}

func CalculatePreV5FileHash(buffer *bytes.Reader, version uint32) (string, error) {
	if version == constants.RecordFileFormatV1 {
		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(buffer)
		if err != nil {
			return "", err
		}
		hash := sha512.Sum384(buf.Bytes())
		return hex.EncodeToString(hash[:]), nil
	} else {
		bytesToRead := make([]byte, constants.PreV5HeaderLength)
		err := binary.Read(buffer, binary.BigEndian, bytesToRead)
		if err != nil {
			return "", err
		}

		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(buffer)
		if err != nil {
			return "", err
		}

		hash := sha512.Sum384(buf.Bytes())
		content := append(bytesToRead, hash[:]...)

		fileHash := sha512.Sum384(content)

		return hex.EncodeToString(fileHash[:]), nil
	}
}

func CalculateV5FileHash(buffer *bytes.Reader) (string, error) {
	buf := new(bytes.Buffer)

	_, err := buf.ReadFrom(buffer)
	if err != nil {
		return "", err
	}

	hash := sha512.Sum384(buf.Bytes())

	return hex.EncodeToString(hash[:]), nil
}

func NewCompactRecordFile(recordFile map[string]interface{}) (*RecordFile, error) {
	valid, err := verifyEndRunningHash(recordFile)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.ErrorInvalidRecordFile
	}

	recordStream, ok := recordFile["record_stream_object"].(string)
	if !ok {
		return nil, errors.ErrorInvalidRecordFile
	}

	record, err := newRecordFromString(recordStream)
	if err != nil {
		return nil, err
	}
	rf := &RecordFile{
		TransactionsMap: make(map[string]*hederaproto.TransactionID),
	}

	err = mapSuccessfulTransactions(rf.TransactionsMap, record.Record)
	if err != nil {
		return nil, err
	}

	metadataHash, err := metadataHash(recordFile)
	if err != nil {
		return nil, err
	}
	rf.MetadataHash = hex.EncodeToString(metadataHash)

	return rf, nil
}

func metadataHash(recordFile map[string]interface{}) ([]byte, error) {
	headStr, ok := recordFile["head"].(string)
	if !ok {
		return nil, errors.ErrorInvalidRecordFile
	}

	h, err := base64.StdEncoding.DecodeString(headStr)
	if err != nil {
		return nil, errors.ErrorInvalidRecordFile
	}

	startRunningHashStr, ok := recordFile["start_running_hash_object"].(string)
	if !ok {
		return nil, errors.ErrorInvalidRecordFile
	}
	srh, err := base64.StdEncoding.DecodeString(startRunningHashStr)
	if err != nil {
		return nil, errors.ErrorInvalidRecordFile
	}

	endRunningHashStr, ok := recordFile["end_running_hash_object"].(string)
	if !ok {
		return nil, errors.ErrorInvalidRecordFile
	}
	erh, err := base64.StdEncoding.DecodeString(endRunningHashStr)
	if err != nil {
		return nil, errors.ErrorInvalidRecordFile
	}

	metadata := make([]byte, 0)
	metadata = append(metadata, h...)
	metadata = append(metadata, srh...)
	metadata = append(metadata, erh...)

	metadataHash := sha512.Sum384(metadata)

	return metadataHash[:], nil
}

func verifyEndRunningHash(recordFile map[string]interface{}) (bool, error) {
	startRunningHashStr, ok := recordFile["start_running_hash_object"].(string)
	if !ok {
		return false, errors.ErrorInvalidRecordFile
	}

	startRunningHash, err := NewHashFromString(startRunningHashStr)
	if err != nil {
		return false, err
	}

	header, err := startRunningHash.Header()
	if err != nil {
		return false, err
	}

	hashesBefore, ok := recordFile["hashes_before"].([]interface{})
	if !ok {
		return false, errors.ErrorInvalidRecordFile
	}

	resultHash := startRunningHash.Hash
	for _, hashBefore := range hashesBefore {
		hashStr, ok := hashBefore.(string)
		if !ok {
			return false, errors.ErrorInvalidRecordFile
		}

		h, err := base64.StdEncoding.DecodeString(hashStr)
		if err != nil {
			return false, errors.ErrorInvalidRecordFile
		}

		resultHash = concatenateAndSha384(header, resultHash, header, h)
	}

	recordStream, ok := recordFile["record_stream_object"].(string)
	if !ok {
		return false, errors.ErrorInvalidRecordFile
	}

	rs, err := base64.StdEncoding.DecodeString(recordStream)
	if err != nil {
		return false, errors.ErrorInvalidRecordFile
	}
	rsHash := sha512.Sum384(rs)
	resultHash = concatenateAndSha384(header, resultHash, header, rsHash[:])

	hashesAfter, ok := recordFile["hashes_after"].([]interface{})
	if !ok {
		return false, errors.ErrorInvalidRecordFile
	}

	for _, hashAfter := range hashesAfter {
		hashStr, ok := hashAfter.(string)
		if !ok {
			return false, errors.ErrorInvalidRecordFile
		}

		h, err := base64.StdEncoding.DecodeString(hashStr)
		if err != nil {
			return false, errors.ErrorInvalidRecordFile
		}

		resultHash = concatenateAndSha384(header, resultHash, header, h)
	}

	endRunningHashStr, ok := recordFile["end_running_hash_object"].(string)
	if !ok {
		return false, errors.ErrorInvalidRecordFile
	}
	endRunningHash, err := NewHashFromString(endRunningHashStr)
	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(endRunningHash.Hash, resultHash), nil
}

func mapSuccessfulTransactions(txMap map[string]*hederaproto.TransactionID, txRecordRawBuffer []byte) error {
	var tr hederaproto.TransactionRecord
	err := proto.Unmarshal(txRecordRawBuffer, &tr)
	if err != nil {
		return err
	}

	transactionReceipt := tr.GetReceipt()
	transactionStatus := transactionReceipt.GetStatus()
	scheduled := tr.ScheduleRef != nil

	if transactionStatus == hederaproto.ResponseCodeEnum_SUCCESS {
		txId := tr.GetTransactionID()
		accId := txId.GetAccountID()
		txTimestamp := txId.GetTransactionValidStart()
		nanos := fmt.Sprintf("%09d", txTimestamp.GetNanos())
		parsedTx := fmt.Sprintf("%d_%d_%d_%d_%s", accId.GetShardNum(), accId.GetRealmNum(), accId.GetAccountNum(), txTimestamp.GetSeconds(), nanos)
		key := fmt.Sprintf("%s-%t", parsedTx, scheduled)
		txMap[key] = txId
	}

	return nil
}

func concatenateAndSha384(prevHeader, previous, currentHeader, current []byte) []byte {
	concat := make([]byte, 0)
	concat = append(concat, prevHeader...)
	concat = append(concat, previous...)
	concat = append(concat, currentHeader...)
	concat = append(concat, current...)
	runningHash := sha512.Sum384(concat)

	return runningHash[:]
}
