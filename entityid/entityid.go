package entityid

import (
	"fmt"
)

const (
	shardBits  int   = 15
	realmBits  int   = 16
	numberBits int   = 32
	shardMask  int64 = (int64(1) << shardBits) - 1
	realmMask  int64 = (int64(1) << realmBits) - 1
	numberMask int64 = (int64(1) << numberBits) - 1
)

// Decode - decodes the Account id into Account string
func Decode(encodedID int64) (string, error) {
	if encodedID < 0 {
		return "", fmt.Errorf("encodedID cannot be negative: %d", encodedID)
	}

	return fmt.Sprintf("%d.%d.%d",
		encodedID>>(realmBits+numberBits), (encodedID>>numberBits)&realmMask, encodedID&numberMask), nil
}
