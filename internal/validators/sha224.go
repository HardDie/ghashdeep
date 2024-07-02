package validators

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
)

type Sha224Validator struct {
}

func NewSha224() Sha224Validator {
	h := Sha224Validator{}
	return h
}

func (v Sha224Validator) Name() string {
	return "sha224"
}

func (v Sha224Validator) Hash(file []byte) []byte {
	hash := sha256.Sum224(file)
	return hash[0:]
}

func (v Sha224Validator) ValidateStream(s io.Reader, hash []byte) (bool, error) {
	h := sha256.New224()
	_, err := io.Copy(h, s)
	if err != nil {
		return false, fmt.Errorf("Sha224Validator.ValidateStream() io.Copy: %w", err)
	}
	fileHash := h.Sum(nil)
	return bytes.Equal(fileHash, hash), nil
}
