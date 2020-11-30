package parser

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"github.com/limechain/hedera-state-proof-verifier-go/internal/errors"
	"github.com/limechain/hedera-state-proof-verifier-go/internal/types"
	"io"
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
	hash := make([]byte, fileHashSize)
	var signature []byte

	for {
		typeDel, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		switch typeDel {
		// hash
		case 4:
			_, err = reader.Read(hash)
			if err != nil {
				return nil, err
			}

			break
		// signature
		case 3:
			sigLength := make([]byte, 4)
			err = binary.Read(reader, binary.BigEndian, sigLength)
			if err != nil {
				return nil, err
			}

			sig := make([]byte, binary.BigEndian.Uint32(sigLength))
			_, err = reader.Read(sig)
			if err != nil {
				return nil, err
			}

			signature = sig
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
