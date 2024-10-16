package validators

import (
	"bytes"
	"crypto/sha512"
	"fmt"
	"io"
)

type Sha384Validator struct {
}

func NewSha384() Sha384Validator {
	h := Sha384Validator{}
	return h
}

func (v Sha384Validator) Name() string {
	return "sha384"
}

func (v Sha384Validator) Hash(file []byte) []byte {
	hash := sha512.Sum384(file)
	return hash[0:]
}

func (v Sha384Validator) CalculateStream(s io.Reader) ([]byte, error) {
	h := sha512.New384()
	_, err := io.Copy(h, s)
	if err != nil {
		return nil, fmt.Errorf("Sha384Validator.CalculateStream() io.Copy: %w", err)
	}
	fileHash := h.Sum(nil)
	return fileHash, nil
}

func (v Sha384Validator) ValidateStream(s io.Reader, hash []byte) (bool, error) {
	h := sha512.New384()
	_, err := io.Copy(h, s)
	if err != nil {
		return false, fmt.Errorf("Sha384Validator.ValidateStream() io.Copy: %w", err)
	}
	fileHash := h.Sum(nil)
	return bytes.Equal(fileHash, hash), nil
}
