package validators

import (
	"bytes"
	"crypto/sha512"
	"fmt"
	"io"
)

type Sha512Validator struct {
}

func NewSha512() Sha512Validator {
	h := Sha512Validator{}
	return h
}

func (v Sha512Validator) Name() string {
	return "sha512"
}

func (v Sha512Validator) Hash(file []byte) []byte {
	hash := sha512.Sum512(file)
	return hash[0:]
}

func (v Sha512Validator) CalculateStream(s io.Reader) ([]byte, error) {
	h := sha512.New()
	_, err := io.Copy(h, s)
	if err != nil {
		return nil, fmt.Errorf("Sha512Validator.CalculateStream() io.Copy: %w", err)
	}
	fileHash := h.Sum(nil)
	return fileHash, nil
}

func (v Sha512Validator) ValidateStream(s io.Reader, hash []byte) (bool, error) {
	h := sha512.New()
	_, err := io.Copy(h, s)
	if err != nil {
		return false, fmt.Errorf("Sha512Validator.ValidateStream() io.Copy: %w", err)
	}
	fileHash := h.Sum(nil)
	return bytes.Equal(fileHash, hash), nil
}
