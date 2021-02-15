package types

import (
	"bytes"
	"encoding/binary"
	"github.com/limechain/hedera-state-proof-verifier-go/internal/constants"
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
