//go:build integration
// +build integration

package crawler

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestIntegration this test verifies that the hash sum calculation
// and verification are compatible and validated.
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
			require.NoErrorf(t, err, "error creating temp dir")
			defer func() {
				err := os.RemoveAll(dir)
				require.NoError(t, err)
			}()

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
			require.NoError(t, err, "error writing data")
			err = file.Close()
			require.NoError(t, err, "error closing file")

			// Calculate hash for created file
			hash := ChooseHashAlg(tc.Validator)
			require.NotNil(t, hash, "error hash not found")
			err = New(hash).Calculate(dir)
			require.NoError(t, err, "error calculating hash")

			err = New(hash).Check(dir)
			require.NoError(t, err, "error checking hash")
		})
	}
}
