package sha384

import (
	"crypto/sha512"
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
	return "sha384"
}

func (v Validator) CalculateStream(s io.Reader) ([]byte, error) {
	h := sha512.New384()
	_, err := io.Copy(h, s)
	if err != nil {
		return nil, fmt.Errorf("sha384.CalculateStream() io.Copy: %w", err)
	}
	fileHash := h.Sum(nil)
	return fileHash, nil
}
