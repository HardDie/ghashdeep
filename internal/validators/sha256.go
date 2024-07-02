package validators

import (
	"bytes"
	"crypto/sha256"
)

type Sha256Validator struct {
}

func NewSha256() Sha256Validator {
	h := Sha256Validator{}
	return h
}

func (v Sha256Validator) Name() string {
	return "sha256"
}

func (v Sha256Validator) Hash(file []byte) []byte {
	hash := sha256.Sum256(file)
	return hash[0:]
}

func (v Sha256Validator) Validate(file, hash []byte) bool {
	fileHash := sha256.Sum256(file)
	return bytes.Equal(fileHash[0:], hash)
}
