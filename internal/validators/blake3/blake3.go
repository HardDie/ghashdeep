package blake3

import (
	"bytes"
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

func (v Validator) Hash(file []byte) []byte {
	hash := blake3.Sum256(file)
	return hash[0:]
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

func (v Validator) ValidateStream(s io.Reader, hash []byte) (bool, error) {
	fileHash, err := v.CalculateStream(s)
	if err != nil {
		return false, fmt.Errorf("blake3.ValidateStream() CalculateStream: %w", err)
	}
	return bytes.Equal(fileHash, hash), nil
}
