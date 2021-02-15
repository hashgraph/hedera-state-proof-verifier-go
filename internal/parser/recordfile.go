package parser

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	hederaproto "github.com/hashgraph/hedera-sdk-go/proto"
	"github.com/limechain/hedera-state-proof-verifier-go/internal/constants"
	"github.com/limechain/hedera-state-proof-verifier-go/internal/errors"
	"github.com/limechain/hedera-state-proof-verifier-go/internal/types"
)

func ParseRecordFile(recordFile string) (map[string]*hederaproto.TransactionID, string, error) {
	bytesRf, err := base64.StdEncoding.DecodeString(recordFile)
	if err != nil {
		return nil, "", err
	}

	bytesReader := bytes.NewReader(bytesRf)
	reader := bufio.NewReader(bytesReader)

	// read record file format version
	intBytes, err := reader.Peek(constants.IntSize)
	if err != nil {
		return nil, "", err
	}
	bytesReader.Reset(bytesRf)

	version := binary.BigEndian.Uint32(intBytes)

	switch version {
	case 1:
	case 2:
		hash, err := types.CalculatePreV5FileHash(bytesReader, version)
		if err != nil {
			return nil, "", err
		}
		bytesReader.Reset(bytesRf)
		mapResult, err := types.NewPreV5RecordFile(bytesReader)
		if err != nil {
			return nil, "", err
		}
		return mapResult, hash, err
	case 5:
		hash, err := types.CalculateV5FileHash(bytesReader)
		if err != nil {
			return nil, "", err
		}
		bytesReader.Reset(bytesRf)
		mapResult, err := types.NewV5RecordFile(bytesReader)
		if err != nil {
			return nil, "", err
		}
		return mapResult, hash, err
	default:
		return nil, "", errors.ErrorUnexpectedTypeDelimiter
	}

	return nil, "", errors.ErrorUnsupportedRecordFileMarker
}
