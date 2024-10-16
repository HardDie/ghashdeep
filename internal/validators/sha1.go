package validators

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
)

type Sha1Validator struct {
}

func NewSha1() Sha1Validator {
	h := Sha1Validator{}
	return h
}

func (v Sha1Validator) Name() string {
	return "sha1"
}

func (v Sha1Validator) Hash(file []byte) []byte {
	hash := sha1.Sum(file)
	return hash[0:]
}

func (v Sha1Validator) CalculateStream(s io.Reader) ([]byte, error) {
	h := sha1.New()
	_, err := io.Copy(h, s)
	if err != nil {
		return nil, fmt.Errorf("Sha1Validator.CalculateStream() io.Copy: %w", err)
	}
	fileHash := h.Sum(nil)
	return fileHash, nil
}

func (v Sha1Validator) ValidateStream(s io.Reader, hash []byte) (bool, error) {
	h := sha1.New()
	_, err := io.Copy(h, s)
	if err != nil {
		return false, fmt.Errorf("Sha1Validator.ValidateStream() io.Copy: %w", err)
	}
	fileHash := h.Sum(nil)
	return bytes.Equal(fileHash, hash), nil
}
