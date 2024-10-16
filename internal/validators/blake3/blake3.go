package blake3

import (
	"fmt"
	"io"

	"github.com/zeebo/blake3"
)

type Validator struct {
}

func New() Validator {
	h := Validator{}
	return h
}

func (v Validator) Name() string {
	return "blake3"
}

func (v Validator) CalculateStream(s io.Reader) ([]byte, error) {
	h := blake3.New()
	_, err := io.Copy(h, s)
	if err != nil {
		return nil, fmt.Errorf("blake3.CalculateStream() io.Copy: %w", err)
	}
	fileHash := h.Sum(nil)
	return fileHash, nil
}
