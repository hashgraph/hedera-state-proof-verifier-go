package parser

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	hederaproto "github.com/hashgraph/hedera-sdk-go/proto"
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
	recordFileWriter := bytes.NewBuffer(result)
	contentsWriter := bytes.NewBuffer(contents)
	index := 0
	version := int(binary.BigEndian.Uint32(bytesRf[index:]))
	recordFileWriter.Write(bytesRf[index : index+4])
	index += 4
	if version >= 2 {
		contentsWriter = bytes.NewBuffer(result)
	}
	recordFileWriter.Write(bytesRf[index : index+4])
	index += 4
	for index < len(bytesRf) {
		typeDel := bytesRf[index]
		index += 1
		switch typeDel {
		case 1:
			recordFileWriter.Write([]byte{typeDel})
			prevHash := bytesRf[index : index+fileHashSize]
			recordFileWriter.Write(prevHash)
			index += fileHashSize
			break
		case 2:
			contentsWriter.Write([]byte{typeDel})
			contentsWriter.Write(bytesRf[index : index+4])
			txRawBytesLength := int(binary.BigEndian.Uint32(bytesRf[index:]))
			index += 4
			contentsWriter.Write(bytesRf[index : index+txRawBytesLength])
			index += txRawBytesLength

			contentsWriter.Write(bytesRf[index : index+4])
			recordRawBytesLength := int(binary.BigEndian.Uint32(bytesRf[index:]))
			index += 4

			transactionRecordRawBuffer := bytesRf[index : index+recordRawBytesLength]
			contentsWriter.Write(transactionRecordRawBuffer)
			index += recordRawBytesLength
			err = parseTransaction(transactionRecordRawBuffer)
			if err != nil {
				return nil, "", err
			}
			break
		default:
			return nil, "", errors.New(fmt.Sprintf(`Unexpected type delimeter %d`, typeDel))
		}
	}

	if version == 2 {
		contentsHash := sha512.Sum384(contentsWriter.Bytes())
		recordFileWriter.Write(contentsHash[:])
	}

	hash := sha512.Sum384(recordFileWriter.Bytes())

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
