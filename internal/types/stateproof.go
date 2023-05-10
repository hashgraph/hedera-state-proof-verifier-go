package types

import (
	"encoding/json"
	"github.com/hashgraph/hedera-state-proof-verifier-go/internal/errors"
)

type StateProof struct {
	RecordFile     interface{}       `json:"record_file"`
	SignatureFiles map[string]string `json:"signature_files"`
	AddressBooks   []string          `json:"address_books"`
	Version        int64             `json:"version"`
}

func NewStateProof(payload []byte) (*StateProof, error) {
	var stateProof StateProof
	err := json.Unmarshal(payload, &stateProof)
	if err != nil {
		return nil, err
	}

	if len(stateProof.SignatureFiles) < 2 {
		return nil, errors.ErrorInvalidSignaturesLength
	}
	if stateProof.RecordFile == nil {
		return nil, errors.ErrorInvalidRecordFile
	}
	if len(stateProof.AddressBooks) < 1 {
		return nil, errors.ErrorInvalidAddressBooksLength
	}

	return &stateProof, nil
}
