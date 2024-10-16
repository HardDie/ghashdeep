package validators

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
)

type Md5Validator struct {
}

func NewMd5() Md5Validator {
	h := Md5Validator{}
	return h
}

func (v Md5Validator) Name() string {
	return "md5"
}

func (v Md5Validator) Hash(file []byte) []byte {
	hash := md5.Sum(file)
	return hash[0:]
}

func (v Md5Validator) CalculateStream(s io.Reader) ([]byte, error) {
	h := md5.New()
	_, err := io.Copy(h, s)
	if err != nil {
		return nil, fmt.Errorf("Md5Validator.CalculateStream() io.Copy: %w", err)
	}
	fileHash := h.Sum(nil)
	return fileHash, nil
}

func (v Md5Validator) ValidateStream(s io.Reader, hash []byte) (bool, error) {
	h := md5.New()
	_, err := io.Copy(h, s)
	if err != nil {
		return false, fmt.Errorf("Md5Validator.ValidateStream() io.Copy: %w", err)
	}
	fileHash := h.Sum(nil)
	return bytes.Equal(fileHash, hash), nil
}
