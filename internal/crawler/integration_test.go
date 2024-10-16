//go:build integration
// +build integration

package crawler

import (
	"os"
	"path"
	"testing"
)

func TestIntegration(t *testing.T) {
	tests := []struct {
		Validator string
	}{
		{Validator: "md5"},
		{Validator: "sha1"},
		{Validator: "sha224"},
		{Validator: "sha256"},
		{Validator: "sha384"},
		{Validator: "sha512"},
		{Validator: "xxhash"},
		{Validator: "blake3"},
	}

	for _, tc := range tests {
		t.Run(tc.Validator, func(t *testing.T) {
			// Create temp dir
			dir, err := os.MkdirTemp("", tc.Validator+"comp")
			if err != nil {
				t.Fatal("error creating temp dir", err)
			}
			defer os.RemoveAll(dir)

			name := "some.txt"
			payload := []byte("hi")

			// Create example file
			filePath := path.Join(dir, name)
			file, err := os.Create(filePath)
			if err != nil {
				// If file can't be created, skip it
				return
			}
			// Write example payload
			_, err = file.Write(payload)
			if err != nil {
				t.Fatal("error writing data", err)
			}
			err = file.Close()
			if err != nil {
				t.Fatal("error closing file", err)
			}

			// Calculate hash for created file
			hash := ChooseHashAlg(tc.Validator)
			if hash == nil {
				t.Fatal("error hash not found")
			}
			err = New(hash).Calculate(dir)
			if err != nil {
				t.Fatal("error calculating hash", err)
			}

			err = New(hash).Check(dir)
			if err != nil {
				t.Fatal("error check hash", err)
			}
		})
	}
}
