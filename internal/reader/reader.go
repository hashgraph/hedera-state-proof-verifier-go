package reader

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/limechain/hedera-state-proof-verifier-go/internal/constants"
)

func LengthAndBytes(buffer *bytes.Reader, minLength, maxLength int, hasChecksum bool) (*uint32, []byte, error) {
	var length uint32
	offset := constants.IntSize
	err := binary.Read(buffer, binary.BigEndian, &length)
	if err != nil {
		return nil, nil, err
	}
	if minLength == maxLength {
		if length != uint32(minLength) {
			return nil, nil, errors.New(fmt.Sprintf("Error reading length and bytes, expect length %d, got %d", minLength, length))
		}
	} else if length < uint32(minLength) || length > uint32(maxLength) {
		return nil, nil, errors.New(fmt.Sprintf("Error reading length and bytes, expect length %d within [%d, %d]", length, minLength, maxLength))
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
			return nil, nil, errors.New(fmt.Sprintf("Error reading length and bytes, expect checksum %d to be %d", checkSum, expected))
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
