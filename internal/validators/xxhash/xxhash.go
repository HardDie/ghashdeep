package xxhash

import (
	"bytes"
	"fmt"
	"io"

	"github.com/cespare/xxhash"
)

type Validator struct {
}

func New() Validator {
	h := Validator{}
	return h
}

func (v Validator) Name() string {
	return "xxhash"
}

func (v Validator) Hash(file []byte) []byte {
	// https://github.com/cespare/xxhash/blob/998dce232f17418a7a5721ecf87ca714025a3243/xxhash.go#L113
	s := xxhash.Sum64(file)
	return append(
		[]byte{},
		byte(s>>56),
		byte(s>>48),
		byte(s>>40),
		byte(s>>32),
		byte(s>>24),
		byte(s>>16),
		byte(s>>8),
		byte(s),
	)
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

func (v Validator) ValidateStream(s io.Reader, hash []byte) (bool, error) {
	fileHash, err := v.CalculateStream(s)
	if err != nil {
		return false, fmt.Errorf("xxhash.ValidateStream() CalculateStream: %w", err)
	}
	return bytes.Equal(fileHash, hash), nil
}
