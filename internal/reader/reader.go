package reader

import (
	"bytes"
	"encoding/binary"
	"github.com/hashgraph/hedera-state-proof-verifier-go/internal/constants"
	"github.com/hashgraph/hedera-state-proof-verifier-go/internal/errors"
)

func LengthAndBytes(buffer *bytes.Reader, minLength, maxLength uint32, hasChecksum bool) (*uint32, []byte, error) {
	var length uint32
	offset := constants.IntSize

	err := binary.Read(buffer, binary.BigEndian, &length)
	if err != nil {
		return nil, nil, err
	}

	if minLength == maxLength {
		if length != minLength {
			return nil, nil, errors.ErrorInvalidLength
		}
	} else if length < minLength || length > maxLength {
		return nil, nil, errors.ErrorInvalidLength
	}

	if hasChecksum {
		var checkSum uint32
		err = binary.Read(buffer, binary.BigEndian, &checkSum)
		if err != nil {
			return nil, nil, err
		}

		expected := constants.SimpleSum - length
		offset += constants.IntSize

		if checkSum != expected {
			return nil, nil, errors.ErrorInvalidChecksum
		}
	}
	finalLength := uint32(offset) + length
	b := make([]byte, length)

	err = binary.Read(buffer, binary.BigEndian, b)
	if err != nil {
		return nil, nil, err
	}

	return &finalLength, b, nil
}
