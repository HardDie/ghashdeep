package xxhash

import (
	"bytes"
	"fmt"
	"io"

	"github.com/cespare/xxhash"
)

type Validator struct {
}

func New() Validator {
	h := Validator{}
	return h
}

func (v Validator) Name() string {
	return "xxhash"
}

func (v Validator) CalculateStream(s io.Reader) ([]byte, error) {
	h := xxhash.New()
	_, err := io.Copy(h, s)
	if err != nil {
		return nil, fmt.Errorf("xxhash.CalculateStream() io.Copy: %w", err)
	}
	fileHash := h.Sum(nil)
	return fileHash, nil
}

func (v Validator) ValidateStream(s io.Reader, hash []byte) (bool, error) {
	fileHash, err := v.CalculateStream(s)
	if err != nil {
		return false, fmt.Errorf("xxhash.ValidateStream() CalculateStream: %w", err)
	}
	return bytes.Equal(fileHash, hash), nil
}
