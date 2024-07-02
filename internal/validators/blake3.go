package validators

import (
	"bytes"
	"fmt"
	"io"

	"github.com/zeebo/blake3"
)

type Blake3Validator struct {
}

func NewBlake3() Blake3Validator {
	h := Blake3Validator{}
	return h
}

func (v Blake3Validator) Name() string {
	return "blake3"
}

func (v Blake3Validator) Hash(file []byte) []byte {
	hash := blake3.Sum256(file)
	return hash[0:]
}

func (v Blake3Validator) ValidateStream(s io.Reader, hash []byte) (bool, error) {
	h := blake3.New()
	_, err := io.Copy(h, s)
	if err != nil {
		return false, fmt.Errorf("Blake3Validator.ValidateStream() io.Copy: %w", err)
	}
	fileHash := h.Sum(nil)
	return bytes.Equal(fileHash[0:], hash), nil
}
