package stateproof

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"github.com/hashgraph/hedera-state-proof-verifier-go/internal/errors"
	"github.com/hashgraph/hedera-state-proof-verifier-go/internal/parser"
	"github.com/hashgraph/hedera-state-proof-verifier-go/internal/types"
	"math"
	"regexp"
)

func Verify(txId string, payload []byte) (bool, error) {
	regex := regexp.MustCompile("[.@\\-]")
	txId = regex.ReplaceAllString(txId, "_")

	stateProof, err := types.NewStateProof(payload)
	if err != nil {
		return false, err
	}

	nodeIdPairs, err := parser.ParseAddressBooks(stateProof.AddressBooks)
	if err != nil {
		return false, err
	}

	signatureFiles, err := parser.ParseSignatureFiles(stateProof.SignatureFiles)
	if err != nil {
		return false, err
	}

	recordFile, err := parser.ParseRecordFile(stateProof.RecordFile)
	if err != nil {
		return false, err
	}
	key := fmt.Sprintf("%s-%t", txId, false)

	if recordFile.TransactionsMap[key] == nil {
		return false, errors.ErrorTransactionNotFound
	}

	err = performStateProof(nodeIdPairs, signatureFiles, recordFile)
	if err != nil {
		return false, err
	}

	return true, nil
}

func performStateProof(nodeIdPubKeyPairs map[string]string, signatureFiles map[string]*types.SignatureFile, recordFile *types.RecordFile) error {
	fileHash, metadataHash, err := verifySignatures(nodeIdPubKeyPairs, signatureFiles)
	if err != nil {
		return err
	}

	if recordFile.Hash != "" && recordFile.Hash == fileHash {
		return nil
	}

	if recordFile.MetadataHash != "" && recordFile.MetadataHash == metadataHash {
		return nil
	}

	return errors.ErrorHashesNotMatch
}

func verifySignatures(nodeIdPubKeyPairs map[string]string, signatureFiles map[string]*types.SignatureFile) (fileHash, metadataHash string, err error) {
	verifiedSigs := make(map[string][]string)
	maxHashCount := 0

	for nodeId, sigFile := range signatureFiles {
		pubKey := nodeIdPubKeyPairs[nodeId]
		if !verifySignature(pubKey, sigFile.Hash, sigFile.Signature) {
			return "", "", errors.ErrorVerifySignature
		}

		if sigFile.MetadataHash != nil && !verifySignature(pubKey, sigFile.MetadataHash, sigFile.MetadataSignature) {
			return "", "", errors.ErrorVerifyMetadataSignature
		}

		hexHash := hex.EncodeToString(sigFile.Hash)
		verifiedSigs[hexHash] = append(verifiedSigs[hexHash], nodeId)

		nodesCount := len(verifiedSigs[hexHash])
		if nodesCount > 1 && nodesCount > maxHashCount {
			maxHashCount = nodesCount
			fileHash = hexHash
			metadataHash = hex.EncodeToString(sigFile.MetadataHash)
		}
	}

	consensusSigs := int(math.Ceil(float64(len(signatureFiles) / 3)))
	if maxHashCount >= consensusSigs {
		return fileHash, metadataHash, nil
	} else {
		return "", "", nil
	}
}

func verifySignature(publicKeyString string, hash []byte, sig []byte) bool {
	pkBytes, err := hex.DecodeString(publicKeyString)
	if err != nil {
		return false
	}
	pk, err := x509.ParsePKIXPublicKey(pkBytes)
	if err != nil {
		return false
	}

	hashedHash := sha512.Sum384(hash)

	res := rsa.VerifyPKCS1v15(pk.(*rsa.PublicKey), crypto.SHA384, hashedHash[:], sig)
	if res != nil {
		return false
	}

	return true
}
