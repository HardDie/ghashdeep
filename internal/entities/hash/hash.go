package hash

import (
	"encoding/hex"
	"fmt"
)

type Hash []byte

func FromString(str string) (Hash, error) {
	bytes, err := hex.DecodeString(str)
	if err != nil {
		return Hash{}, fmt.Errorf("hex.DecodeString(%s): %w", str, err)
	}
	return bytes, nil
}

func FromBytes(b []byte) Hash {
	return b
}

func (h Hash) String() string {
	return hex.EncodeToString([]byte(h))
}
