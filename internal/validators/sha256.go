package validators

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
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

func (v Sha256Validator) ValidateStream(s io.Reader, hash []byte) (bool, error) {
	h := sha256.New()
	_, err := io.Copy(h, s)
	if err != nil {
		return false, fmt.Errorf("XxhashValidator.ValidateStream() io.Copy: %w", err)
	}
	fileHash := h.Sum(nil)
	return bytes.Equal(fileHash, hash), nil
}
