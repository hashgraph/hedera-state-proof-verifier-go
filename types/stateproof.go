package types

import (
	"encoding/json"
	"github.com/limechain/hedera-state-proof-verifier-go/errors"
)

type StateProof struct {
	RecordFile     string            `json:"record_file"`
	SignatureFiles map[string]string `json:"signature_files"`
	AddressBooks   []string          `json:"address_books"`
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
	if stateProof.RecordFile == "" {
		return nil, errors.ErrorInvalidRecordFile
	}
	if len(stateProof.AddressBooks) < 1 {
		return nil, errors.ErrorInvalidAddressBooksLength
	}

	return &stateProof, nil
}
