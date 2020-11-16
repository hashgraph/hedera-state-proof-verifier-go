package parser

import (
	"encoding/base64"
	"encoding/binary"
	"github.com/limechain/hedera-state-proof-verifier-go/errors"
	"github.com/limechain/hedera-state-proof-verifier-go/types"
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
	index := 0
	var hash []byte
	var signature []byte
	for index < len(bytesSigFile) {
		typeDel := bytesSigFile[index]
		index += 1

		switch typeDel {
		// hash
		case 4:
			hash = bytesSigFile[index : index+fileHashSize]
			index += fileHashSize
			break
		// signature
		case 3:
			signatureLength := int(binary.BigEndian.Uint32(bytesSigFile[index:]))
			index += 4
			signature = bytesSigFile[index : index+signatureLength]
			index += signatureLength
			break
		default:
			return nil, errors.ErrorUnexpectedTypeDelimiter
		}
	}
	return &types.SignatureFile{
		Hash:      hash,
		Signature: signature,
	}, nil
}
