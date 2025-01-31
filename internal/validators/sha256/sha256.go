package sha256

import (
	"crypto/sha256"
	"fmt"
	"io"
)

type Validator struct{}

func New() Validator {
	h := Validator{}
	return h
}

func (v Validator) Name() string {
	return "sha256"
}

func (v Validator) CalculateStream(s io.Reader) ([]byte, error) {
	h := sha256.New()
	_, err := io.Copy(h, s)
	if err != nil {
		return nil, fmt.Errorf("sha256.CalculateStream() io.Copy: %w", err)
	}
	fileHash := h.Sum(nil)
	return fileHash, nil
}
