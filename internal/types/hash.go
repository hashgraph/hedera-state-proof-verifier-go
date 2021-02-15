package types

import (
	"bytes"
	"encoding/binary"
	"github.com/limechain/hedera-state-proof-verifier-go/internal/constants"
	"github.com/limechain/hedera-state-proof-verifier-go/internal/reader"
)

type HashStruct struct {
	Stream
	BodyLength uint32
	DigestType []byte
	Hash       []byte
}

func NewHashStruct(buffer *bytes.Reader) (*HashStruct, error) {
	stream, err := NewStream(buffer)
	if err != nil {
		return nil, err
	}
	hashStruct := &HashStruct{
		Stream: *stream,
	}

	bodyLength, err := hashStruct.readBody(buffer)
	if err != nil {
		return nil, err
	}

	hashStruct.BodyLength = *bodyLength + stream.BodyLength

	return hashStruct, nil
}

func (h *HashStruct) readBody(buffer *bytes.Reader) (*uint32, error) {
	h.DigestType = make([]byte, constants.IntSize)
	err := binary.Read(buffer, binary.BigEndian, h.DigestType)
	if err != nil {
		return nil, err
	}

	length, b, err := reader.LengthAndBytes(buffer, constants.Sha384Length, constants.Sha384Length, false)
	if err != nil {
		return nil, err
	}
	h.Hash = b

	finalLength := *length + constants.IntSize

	return &finalLength, nil
}
