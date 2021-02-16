package parser

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"github.com/limechain/hedera-state-proof-verifier-go/internal/constants"
	"github.com/limechain/hedera-state-proof-verifier-go/internal/errors"
	"github.com/limechain/hedera-state-proof-verifier-go/internal/types"
)

func ParseRecordFile(record string) (*types.RecordFile, error) {
	bytesRf, err := base64.StdEncoding.DecodeString(record)
	if err != nil {
		return nil, err
	}

	bytesReader := bytes.NewReader(bytesRf)
	reader := bufio.NewReader(bytesReader)

	// record file format version
	intBytes, err := reader.Peek(constants.IntSize)
	if err != nil {
		return nil, err
	}
	bytesReader.Reset(bytesRf)

	version := binary.BigEndian.Uint32(intBytes)

	var recordFile *types.RecordFile
	switch version {
	case constants.RecordFileFormatV1:
	case constants.RecordFileFormatV2:
		hash, err := types.CalculatePreV5FileHash(bytesReader, version)
		if err != nil {
			return nil, err
		}
		bytesReader.Reset(bytesRf)

		recordFile, err = types.NewPreV5RecordFile(bytesReader)
		if err != nil {
			return nil, err
		}

		recordFile.Hash = hash
		break
	case constants.RecordFileFormatV5:
		hash, err := types.CalculateV5FileHash(bytesReader)
		if err != nil {
			return nil, err
		}

		bytesReader.Reset(bytesRf)

		recordFile, err = types.NewV5RecordFile(bytesReader)
		if err != nil {
			return nil, err
		}
		recordFile.Hash = hash

		break
	default:
		return nil, errors.ErrorUnexpectedTypeDelimiter
	}

	return recordFile, nil
}
