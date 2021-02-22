package parser

import (
	"bytes"
	"encoding/base64"
	"github.com/hashgraph/hedera-state-proof-verifier-go/internal/constants"
	"github.com/hashgraph/hedera-state-proof-verifier-go/internal/errors"
	"github.com/hashgraph/hedera-state-proof-verifier-go/internal/types"
)

func ParseSignatureFiles(signatureFiles map[string]string) (map[string]*types.SignatureFile, error) {
	result := make(map[string]*types.SignatureFile)

	for key, signatureFile := range signatureFiles {
		bytes, err := base64.StdEncoding.DecodeString(signatureFile)
		if err != nil {
			return nil, err
		}

		signatureFile, err := parseSignatureFile(bytes)
		if err != nil {
			return nil, err
		}
		result[key] = signatureFile
	}

	return result, nil
}

func parseSignatureFile(bytesSigFile []byte) (*types.SignatureFile, error) {
	reader := bytes.NewReader(bytesSigFile)

	fileFormatVersion, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	switch fileFormatVersion {
	case constants.SignatureFileFormatV4:
		return types.NewV2SignatureFile(reader)
	case constants.SignatureFileFormatV5:
		return types.NewV5SignatureFile(reader)
	default:
		return nil, errors.ErrorUnexpectedTypeDelimiter
	}
}
