package types

import (
	"bytes"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/golang/protobuf/proto"
	hederaproto "github.com/hashgraph/hedera-sdk-go/proto"
	"github.com/limechain/hedera-state-proof-verifier-go/internal/constants"
	"github.com/limechain/hedera-state-proof-verifier-go/internal/errors"
	"github.com/limechain/hedera-state-proof-verifier-go/internal/reader"
	"io"
)

type record struct {
	Stream
	BodyLength  uint32
	Record      []byte
	Transaction []byte
}

type RecordFile struct {
	Hash            string
	TransactionsMap map[string]*hederaproto.TransactionID
}

func newRecord(buffer *bytes.Reader) (*record, error) {
	stream, err := NewStream(buffer)
	if err != nil {
		return nil, err
	}
	recordFile := &record{
		Stream: *stream,
	}

	bodyLength, err := recordFile.readBody(buffer)
	if err != nil {
		return nil, err
	}

	recordFile.BodyLength = *bodyLength + stream.BodyLength

	return recordFile, nil
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

func CalculatePreV5FileHash(buffer *bytes.Reader, version uint32) (string, error) {
	if version == 1 {
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

func mapSuccessfulTransactions(txMap map[string]*hederaproto.TransactionID, txRecordRawBuffer []byte) error {
	var tr hederaproto.TransactionRecord
	err := proto.Unmarshal(txRecordRawBuffer, &tr)
	if err != nil {
		return err
	}

	transactionReceipt := tr.GetReceipt()
	transactionStatus := transactionReceipt.GetStatus()

	if transactionStatus == hederaproto.ResponseCodeEnum_SUCCESS {
		txId := tr.GetTransactionID()
		accId := txId.GetAccountID()
		txTimestamp := txId.GetTransactionValidStart()
		nanos := fmt.Sprintf("%09d", txTimestamp.GetNanos())
		parsedTx := fmt.Sprintf("%d_%d_%d_%d_%s", accId.GetShardNum(), accId.GetRealmNum(), accId.GetAccountNum(), txTimestamp.GetSeconds(), nanos)
		txMap[parsedTx] = txId
	}

	return nil
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

	return recordFile, nil
}
