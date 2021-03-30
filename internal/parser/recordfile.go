package parser

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"github.com/hashgraph/hedera-state-proof-verifier-go/internal/constants"
	"github.com/hashgraph/hedera-state-proof-verifier-go/internal/errors"
	"github.com/hashgraph/hedera-state-proof-verifier-go/internal/types"
)

func ParseRecordFile(record interface{}) (*types.RecordFile, error) {
	recordFileType, err := checkType(record)
	if err != nil {
		return nil, err
	}

	switch recordFileType {
	case constants.FullRecordFile:
		return parseFullRecordFile(record.(string))
	case constants.CompactRecordFile:
		return parseCompactRecordFile(record.(map[string]interface{}))
	default:
		return nil, errors.ErrorInvalidRecordFile
	}
}

func parseCompactRecordFile(record map[string]interface{}) (*types.RecordFile, error) {
	return types.NewCompactRecordFile(record)
}

func parseFullRecordFile(record string) (*types.RecordFile, error) {
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
	default:
		return nil, errors.ErrorUnexpectedTypeDelimiter
	}

	return recordFile, nil
}

func checkType(recordFile interface{}) (recordFileType constants.RecordFileType, error error) {
	switch recordFile.(type) {
	case string:
		version, err := readVersion(recordFile.(string))
		if err != nil {
			return constants.InvalidRecordFile, err
		}

		validType := version == constants.RecordFileFormatV1 ||
			version == constants.RecordFileFormatV2 ||
			version == constants.RecordFileFormatV5
		if validType {
			return constants.FullRecordFile, nil
		}
		return constants.InvalidRecordFile, nil
	case map[string]interface{}:
		record := recordFile.(map[string]interface{})
		version, err := readVersion(record["head"].(string))
		if err != nil {
			return constants.InvalidRecordFile, err
		}

		validType := version == constants.RecordFileFormatV5
		if validType {
			return constants.CompactRecordFile, nil
		}
		return constants.InvalidRecordFile, nil
	default:
		return constants.InvalidRecordFile, nil
	}
}

func readVersion(s string) (uint32, error) {
	bytesRf, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return 0, err
	}

	bytesReader := bytes.NewReader(bytesRf)
	reader := bufio.NewReader(bytesReader)
	intBytes, err := reader.Peek(constants.IntSize)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint32(intBytes), nil
}
