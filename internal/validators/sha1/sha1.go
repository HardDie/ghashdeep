package sha1

import (
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

func (v Validator) CalculateStream(s io.Reader) ([]byte, error) {
	h := sha1.New()
	_, err := io.Copy(h, s)
	if err != nil {
		return nil, fmt.Errorf("sha1.CalculateStream() io.Copy: %w", err)
	}
	fileHash := h.Sum(nil)
	return fileHash, nil
}
