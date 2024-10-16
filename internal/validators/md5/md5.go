package md5

import (
	"bytes"
	"crypto/md5"
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
	return "md5"
}

func (v Validator) Hash(file []byte) []byte {
	hash := md5.Sum(file)
	return hash[0:]
}

func (v Validator) CalculateStream(s io.Reader) ([]byte, error) {
	h := md5.New()
	_, err := io.Copy(h, s)
	if err != nil {
		return nil, fmt.Errorf("md5.CalculateStream() io.Copy: %w", err)
	}
	fileHash := h.Sum(nil)
	return fileHash, nil
}

func (v Validator) ValidateStream(s io.Reader, hash []byte) (bool, error) {
	fileHash, err := v.CalculateStream(s)
	if err != nil {
		return false, fmt.Errorf("md5.ValidateStream() CalculateStream: %w", err)
	}
	return bytes.Equal(fileHash, hash), nil
}
