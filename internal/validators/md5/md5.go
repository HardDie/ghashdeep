package md5

import (
	//nolint:gosec //not for security
	"crypto/md5"
	"fmt"
	"io"
)

type Validator struct{}

func New() Validator {
	h := Validator{}
	return h
}

func (v Validator) Name() string {
	return "md5"
}

func (v Validator) CalculateStream(s io.Reader) ([]byte, error) {
	//nolint:gosec //not for security
	h := md5.New()
	_, err := io.Copy(h, s)
	if err != nil {
		return nil, fmt.Errorf("md5.CalculateStream() io.Copy: %w", err)
	}
	fileHash := h.Sum(nil)
	return fileHash, nil
}
