package parser

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/golang/protobuf/proto"
	hederaproto "github.com/hashgraph/hedera-sdk-go/proto"
	"github.com/limechain/hedera-state-proof-verifier-go/internal/errors"
	"io"
)

const (
	fileHashSize = 48
)

var (
	txMap = make(map[string]*hederaproto.TransactionID)
)

func ParseRecordFile(recordFile string) (map[string]*hederaproto.TransactionID, string, error) {
	bytesRf, err := base64.StdEncoding.DecodeString(recordFile)
	if err != nil {
		return nil, "", err
	}

	var result []byte
	var contents []byte
	reader := bytes.NewReader(bytesRf)
	recordFileWriter := bytes.NewBuffer(result)
	contentsWriter := bytes.NewBuffer(contents)

	intBytes := make([]byte, 4)
	// read record file format version
	err = binary.Read(reader, binary.BigEndian, intBytes)
	if err != nil {
		return nil, "", err
	}
	recordFileWriter.Write(intBytes)

	version := binary.BigEndian.Uint32(intBytes)
	if version >= 2 {
		contentsWriter = bytes.NewBuffer(result)
	}

	// version
	err = binary.Read(reader, binary.BigEndian, intBytes)
	if err != nil {
		return nil, "", err
	}
	recordFileWriter.Write(intBytes)

	for {
		typeDel, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, "", err
		}

		switch typeDel {
		// RECORD_TYPE_PREV_HASH
		case 1:
			recordFileWriter.Write([]byte{typeDel})
			fileHashBytes := make([]byte, fileHashSize)
			_, err = reader.Read(fileHashBytes)
			if err != nil {
				return nil, "", err
			}

			recordFileWriter.Write(fileHashBytes)
			break
		// RECORD_TYPE_RECORD
		case 2:
			contentsWriter.Write([]byte{typeDel})

			// transaction raw bytes
			err = binary.Read(reader, binary.BigEndian, intBytes)
			if err != nil {
				return nil, "", err
			}
			contentsWriter.Write(intBytes)

			txRawLength := binary.BigEndian.Uint32(intBytes)
			txRaw := make([]byte, txRawLength)
			_, err = reader.Read(txRaw)
			if err != nil {
				return nil, "", err
			}
			contentsWriter.Write(txRaw)

			// tx record raw bytes
			err = binary.Read(reader, binary.BigEndian, intBytes)
			if err != nil {
				return nil, "", err
			}
			contentsWriter.Write(intBytes)

			recordRawLength := binary.BigEndian.Uint32(intBytes)

			// tx record raw buffer
			txRecordRaw := make([]byte, recordRawLength)
			_, err = reader.Read(txRecordRaw)
			if err != nil {
				return nil, "", err
			}
			contentsWriter.Write(txRecordRaw)

			err = parseTransaction(txRecordRaw)
			if err != nil {
				return nil, "", err
			}
			break
		default:
			return nil, "", errors.ErrorUnexpectedTypeDelimiter
		}
	}

	if version == 2 {
		contentsHash := sha512.Sum384(contentsWriter.Bytes())
		recordFileWriter.Write(contentsHash[:])
	}

	hash := sha512.Sum384(recordFileWriter.Bytes())

	// set record file hash
	recordFileHash := hex.EncodeToString(hash[:])

	return txMap, recordFileHash, nil
}

func parseTransaction(txRecordRawBuffer []byte) error {
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
