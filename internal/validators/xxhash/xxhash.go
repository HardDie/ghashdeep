package xxhash

import (
	"fmt"
	"io"

	"github.com/cespare/xxhash"
)

type Validator struct{}

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
