package types

import (
	"bytes"
	"encoding/binary"
	"github.com/limechain/hedera-state-proof-verifier-go/internal/constants"
	"github.com/limechain/hedera-state-proof-verifier-go/internal/reader"
)

type Hash struct {
	Stream
	BodyLength uint32
	DigestType []byte
	Hash       []byte
}

func NewHash(buffer *bytes.Reader) (*Hash, error) {
	stream, err := NewStream(buffer)
	if err != nil {
		return nil, err
	}

	hash := &Hash{
		Stream: *stream,
	}

	bodyLength, err := hash.readBody(buffer)
	if err != nil {
		return nil, err
	}

	hash.BodyLength = *bodyLength + stream.BodyLength

	return hash, nil
}

func (h *Hash) readBody(buffer *bytes.Reader) (*uint32, error) {
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
