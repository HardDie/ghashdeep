package validators

import "encoding/hex"

func calcHashLen(
	hash interface {
		Hash(file []byte) []byte
	},
) int {
	return len(hex.EncodeToString(hash.Hash([]byte{1})))
}
