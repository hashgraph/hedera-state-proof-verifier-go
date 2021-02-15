package types

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/limechain/hedera-state-proof-verifier-go/internal/constants"
	"github.com/limechain/hedera-state-proof-verifier-go/internal/reader"
)

type SignatureFile struct {
	Stream
	Hash              []byte
	Signature         []byte
	Version           int
	MetadataHash      []byte
	MetadataSignature []byte
}

func NewSignatureFile(buffer *bytes.Reader) (*SignatureFile, error) {
	stream, err := NewStream(buffer)
	if err != nil {
		return nil, err
	}
	signatureFile := &SignatureFile{
		Stream: *stream,
	}

	bodyLength, err := signatureFile.readBody(buffer)
	if err != nil {
		return nil, err
	}

	signatureFile.BodyLength = *bodyLength

	return signatureFile, nil
}

func (sf *SignatureFile) readBody(buffer *bytes.Reader) (*uint32, error) {
	var sigFileType uint32
	err := binary.Read(buffer, binary.BigEndian, &sigFileType)
	if err != nil {
		return nil, err
	}

	if sigFileType != constants.Sha384WithRsaType {
		return nil, errors.New(fmt.Sprintf("Error reading Signature object, expected type %d, got %d", constants.Sha384WithRsaType, sigFileType))
	}

	length, b, err := reader.LengthAndBytes(buffer, constants.ByteSize, constants.Sha384WithRsaMaxLength, true)
	if err != nil {
		return nil, err
	}
	sf.Signature = b

	finalLength := *length + constants.IntSize

	return &finalLength, nil
}

func NewV2SignatureFile(buffer *bytes.Reader) (*SignatureFile, error) {
	hash := make([]byte, constants.Sha384Length)

	_, err := buffer.Read(hash)
	if err != nil {
		return nil, err
	}

	typeDel, err := buffer.ReadByte()
	if err != nil {
		return nil, err
	}

	if typeDel != 3 {

	}

	_, b, err := reader.LengthAndBytes(buffer, constants.ByteSize, constants.Sha384WithRsaMaxLength, false)
	if err != nil {
		return nil, err
	}

	if buffer.Len() != 0 {
		return nil, errors.New("extra data discovered in signature file")
	}

	return &SignatureFile{
		Stream:    Stream{},
		Hash:      hash,
		Signature: b,
		Version:   2,
	}, nil
}

func NewV5SignatureFile(buffer *bytes.Reader) (*SignatureFile, error) {
	err := binary.Read(buffer, binary.BigEndian, make([]byte, constants.IntSize))
	if err != nil {
		return nil, err
	}

	hashStruct, err := NewHashStruct(buffer)
	if err != nil {
		return nil, err
	}

	sigFile, err := NewSignatureFile(buffer)
	if err != nil {
		return nil, err
	}

	metadataHash, err := NewHashStruct(buffer)
	if err != nil {
		return nil, err
	}

	metadataSigFile, err := NewSignatureFile(buffer)
	if err != nil {
		return nil, err
	}

	if buffer.Len() != 0 {
		return nil, errors.New("extra data discovered in signature file")
	}

	return &SignatureFile{
		Stream:            Stream{},
		Hash:              hashStruct.Hash,
		Signature:         sigFile.Signature,
		Version:           5,
		MetadataHash:      metadataHash.Hash,
		MetadataSignature: metadataSigFile.Signature,
	}, nil
}
