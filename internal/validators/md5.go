package validators

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
)

type Md5Validator struct {
	len int
}

func NewMd5() Md5Validator {
	h := Md5Validator{}
	h.len = calcHashLen(h)
	return h
}

func (v Md5Validator) Name() string {
	return "md5"
}

func (v Md5Validator) Len() int {
	return v.len
}

func (v Md5Validator) Hash(file []byte) []byte {
	hash := md5.Sum(file)
	return hash[0:]
}

func (v Md5Validator) Validate(file, hash []byte) bool {
	fileHash := md5.Sum(file)
	return bytes.Equal(fileHash[0:], hash)
}

func (v Md5Validator) ValidateStream(s io.Reader, hash []byte) (bool, error) {
	h := md5.New()
	_, err := io.Copy(h, s)
	if err != nil {
		return false, fmt.Errorf("Md5Validator.ValidateStream() io.Copy: %w", err)
	}
	fileHash := h.Sum(nil)
	return bytes.Equal(fileHash[0:], hash), nil
}
