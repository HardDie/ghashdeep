//go:build integration

package crawler

import (
	"strings"
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
			dir := CreateTempDir(t)
			CreateDummyFile(t, dir, "some.txt", []byte("hi"))

			AppCalculateHash(t, tc.Validator, dir)
			AppCheckHash(t, tc.Validator, dir)
		})
	}
}

// TestHashValue calculate the hash using the library method
// and parse stdout from the system binary using the exec call and compare the hash values.
func TestHashValue(t *testing.T) {
	tests := []struct {
		Validator string
		AppName   string
	}{
		{Validator: "md5", AppName: "md5sum"},
		{Validator: "sha1", AppName: "sha1sum"},
		{Validator: "sha224", AppName: "sha224sum"},
		{Validator: "sha256", AppName: "sha256sum"},
		{Validator: "sha384", AppName: "sha384sum"},
		{Validator: "sha512", AppName: "sha512sum"},
		{Validator: "xxhash", AppName: "xxh64sum"},
		{Validator: "blake3", AppName: "b3sum"},
	}

	data := []byte("hi")

	for _, tc := range tests {
		t.Run(tc.Validator, func(t *testing.T) {
			if err := CheckCmdExists(tc.AppName); err != nil {
				t.Skip(err)
			}

			hashHex := AppCalculateHashString(t, tc.Validator, data)
			stdout := CmdCalculateHastStdout(t, tc.AppName, data)

			resp := strings.Split(stdout, " ")
			require.Equal(t, hashHex, resp[0])
		})
	}
}

// TestAppCalculateCmdCheck calculate the hash using the library
// and verify the output file with the binary from the system using an exec call.
func TestAppCalculateCmdCheck(t *testing.T) {
	tests := []struct {
		Validator   string
		AppName     string
		CustomFlags []string
	}{
		{Validator: "md5", AppName: "md5sum"},
		{Validator: "sha1", AppName: "sha1sum"},
		{Validator: "sha224", AppName: "sha224sum"},
		{Validator: "sha256", AppName: "sha256sum"},
		{Validator: "sha384", AppName: "sha384sum"},
		{Validator: "sha512", AppName: "sha512sum"},
		{Validator: "xxhash", AppName: "xxh64sum"},
		{
			Validator:   "blake3",
			AppName:     "b3sum",
			CustomFlags: []string{""},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Validator, func(t *testing.T) {
			if err := CheckCmdExists(tc.AppName); err != nil {
				t.Skip(err)
			}

			dir := CreateTempDir(t)
			CreateDummyFile(t, dir, "some.txt", []byte("hi"))

			AppCalculateHash(t, tc.Validator, dir)
			CmdCheckHash(t, tc.AppName, tc.Validator, dir, tc.CustomFlags)
		})
	}
}

// TestCmdCalculateAppCheck calculate the hash using the binary from the system
// using the exec call, and validate the output file using the library method.
func TestCmdCalculateAppCheck(t *testing.T) {
	tests := []struct {
		Validator string
		AppName   string
	}{
		{Validator: "md5", AppName: "md5sum"},
		{Validator: "sha1", AppName: "sha1sum"},
		{Validator: "sha224", AppName: "sha224sum"},
		{Validator: "sha256", AppName: "sha256sum"},
		{Validator: "sha384", AppName: "sha384sum"},
		{Validator: "sha512", AppName: "sha512sum"},
		{Validator: "xxhash", AppName: "xxh64sum"},
		{Validator: "blake3", AppName: "b3sum"},
	}

	for _, tc := range tests {
		t.Run(tc.Validator, func(t *testing.T) {
			if err := CheckCmdExists(tc.AppName); err != nil {
				t.Skip(err)
			}

			dir := CreateTempDir(t)
			CreateDummyFile(t, dir, "some.txt", []byte("hi"))

			CmdCalculateHash(t, tc.AppName, tc.Validator, dir)
			AppCheckHash(t, tc.Validator, dir)
		})
	}
}
