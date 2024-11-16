package sha224

import (
	"bytes"
	"encoding/hex"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckHash(t *testing.T) {
	payload := []byte("hi")
	appName := "sha224sum"

	// Check if hash app is installed
	_, err := exec.LookPath(appName)
	if err != nil {
		t.Skip(err)
	}

	// Calculating hash using internal library
	validator := New()
	hash, err := validator.CalculateStream(bytes.NewReader(payload))
	require.NoErrorf(t, err, "error calculating hash by library")
	hashHex := hex.EncodeToString(hash)

	// Calculating hash using system app
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd := exec.Command(appName)
	cmd.Stdin = bytes.NewReader(payload)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err = cmd.Run()
	require.NoErrorf(t, err, "error calculating hash by application: %s", stderr.String())

	// Get hash hex value from output
	resp := strings.Split(stdout.String(), " ")

	// Compare result
	require.Equal(t, hashHex, resp[0])
}
