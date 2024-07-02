package validators

import (
	"bytes"
	"fmt"
	"io"

	"github.com/cespare/xxhash"
)

type XxhashValidator struct {
}

func NewXxhash() XxhashValidator {
	h := XxhashValidator{}
	return h
}

func (v XxhashValidator) Name() string {
	return "xxhash"
}

func (v XxhashValidator) Hash(file []byte) []byte {
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

func (v XxhashValidator) ValidateStream(s io.Reader, hash []byte) (bool, error) {
	h := xxhash.New()
	_, err := io.Copy(h, s)
	if err != nil {
		return false, fmt.Errorf("XxhashValidator.ValidateStream() io.Copy: %w", err)
	}
	fileHash := h.Sum(nil)
	return bytes.Equal(fileHash, hash), nil
}
