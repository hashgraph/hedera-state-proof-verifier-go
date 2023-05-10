package types

import (
	"bytes"
	"encoding/binary"
	"github.com/hashgraph/hedera-state-proof-verifier-go/internal/constants"
)

type Stream struct {
	ClassId      int64
	ClassVersion int32
	BodyLength   uint32
}

func NewStream(buffer *bytes.Reader) (*Stream, error) {
	stream := &Stream{}
	err := binary.Read(buffer, binary.BigEndian, &stream.ClassId)
	if err != nil {
		return nil, err
	}

	err = binary.Read(buffer, binary.BigEndian, &stream.ClassVersion)
	if err != nil {
		return nil, err
	}

	stream.BodyLength = constants.IntSize + constants.LongSize

	return stream, nil
}

func (s *Stream) Header() ([]byte, error) {
	buffer := bytes.NewBuffer(make([]byte, 0))

	err := binary.Write(buffer, binary.LittleEndian, s.ClassId)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.LittleEndian, s.ClassVersion)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
