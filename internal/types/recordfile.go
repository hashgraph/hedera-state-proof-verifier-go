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

var (
	txMap = make(map[string]*hederaproto.TransactionID)
)

type RecordFileStruct struct {
	Stream
	BodyLength  uint32
	Record      []byte
	Transaction []byte
}

func NewRecordFileStruct(buffer *bytes.Reader) (*RecordFileStruct, error) {
	stream, err := NewStream(buffer)
	if err != nil {
		return nil, err
	}
	recordFileStruct := &RecordFileStruct{
		Stream: *stream,
	}

	bodyLength, err := recordFileStruct.readBody(buffer)
	if err != nil {
		return nil, err
	}

	recordFileStruct.BodyLength = *bodyLength + stream.BodyLength

	return recordFileStruct, nil
}

func (rfs *RecordFileStruct) readBody(buffer *bytes.Reader) (*uint32, error) {
	recordLength, recordBytes, err := reader.LengthAndBytes(buffer, constants.ByteSize, constants.MaxRecordLength, false)
	if err != nil {
		return nil, err
	}

	transactionLength, transactionBytes, err := reader.LengthAndBytes(buffer, constants.ByteSize, constants.MaxTransactionLength, false)
	if err != nil {
		return nil, err
	}

	rfs.Record = recordBytes
	rfs.Transaction = transactionBytes

	totalLength := *recordLength + *transactionLength

	return &totalLength, nil
}

func NewPreV5RecordFile(buffer *bytes.Reader) (map[string]*hederaproto.TransactionID, error) {
	bytesToRead := make([]byte, constants.PreV5HeaderLength)
	err := binary.Read(buffer, binary.BigEndian, bytesToRead)
	if err != nil {
		return nil, err
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

		err = mapSuccessfulTransactions(recordBytes)
		if err != nil {
			return nil, err
		}
	}

	return txMap, nil
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
		concat := append(bytesToRead, hash[:]...)

		res := sha512.Sum384(concat)
		return hex.EncodeToString(res[:]), nil
	}
}

func mapSuccessfulTransactions(txRecordRawBuffer []byte) error {
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

func NewV5RecordFile(buffer *bytes.Reader) (map[string]*hederaproto.TransactionID, error) {
	var metadata []byte
	hashOffset := make([]byte, constants.V5StartHashOffset)
	err := binary.Read(buffer, binary.BigEndian, hashOffset)
	if err != nil {
		return nil, err
	}
	metadata = append(metadata, hashOffset...)

	hashStruct, err := NewHashStruct(buffer)
	if err != nil {
		return nil, err
	}

	_, err = buffer.Seek(-int64(hashStruct.BodyLength), io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	hashStructBytes := make([]byte, hashStruct.BodyLength)
	err = binary.Read(buffer, binary.BigEndian, hashStructBytes)
	if err != nil {
		return nil, err
	}
	metadata = append(metadata, hashStructBytes...)

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

		if classId == hashStruct.ClassId {
			break
		}

		recordStruct, err := NewRecordFileStruct(buffer)
		if err != nil {
			return nil, err
		}

		err = mapSuccessfulTransactions(recordStruct.Record)
		if err != nil {
			return nil, err
		}
	}
	finalHashStruct, err := NewHashStruct(buffer)
	if err != nil {
		return nil, err
	}

	if buffer.Len() != 0 {
		return nil, errors.ErrorExtraDataInRecordFile
	}

	_, err = buffer.Seek(-int64(finalHashStruct.BodyLength), io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buffer, binary.BigEndian, hashStructBytes)
	if err != nil {
		return nil, err
	}
	metadata = append(metadata, hashStructBytes...)

	return txMap, nil
}
