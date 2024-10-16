package sha1

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
)

type Validator struct {
}

func New() Validator {
	h := Validator{}
	return h
}

func (v Validator) Name() string {
	return "sha1"
}

func (v Validator) Hash(file []byte) []byte {
	hash := sha1.Sum(file)
	return hash[0:]
}

func (v Validator) CalculateStream(s io.Reader) ([]byte, error) {
	h := sha1.New()
	_, err := io.Copy(h, s)
	if err != nil {
		return nil, fmt.Errorf("sha1.CalculateStream() io.Copy: %w", err)
	}
	fileHash := h.Sum(nil)
	return fileHash, nil
}

func (v Validator) ValidateStream(s io.Reader, hash []byte) (bool, error) {
	fileHash, err := v.CalculateStream(s)
	if err != nil {
		return false, fmt.Errorf("sha1.ValidateStream() CalculateStream: %w", err)
	}
	return bytes.Equal(fileHash, hash), nil
}
