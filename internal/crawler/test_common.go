package crawler

import (
	"bytes"
	"encoding/hex"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/HardDie/ghashdeep/internal/entities/config"
)

func CreateTempDir(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

func CreateDummyFile(t *testing.T, dir, name string, data []byte) {
	t.Helper()
	filePath := path.Join(dir, name)
	file, err := os.Create(filePath)
	require.NoErrorf(t, err, "error creating file: %s", filePath)
	_, err = file.Write(data)
	require.NoError(t, err, "error writing data")
	err = file.Close()
	require.NoError(t, err, "error closing file")
}

func CheckCmdExists(appName string) error {
	_, err := exec.LookPath(appName)
	if err != nil {
		return err
	}
	return nil
}

func AppCalculateHashString(t *testing.T, hashAlg string, data []byte) string {
	t.Helper()
	hash := ChooseHashAlg(hashAlg)
	require.NotNil(t, hash, "error hash not found")
	resp, err := hash.CalculateStream(bytes.NewReader(data))
	require.NoErrorf(t, err, "error calculating hash by library")
	return hex.EncodeToString(resp)
}

func AppCalculateHash(t *testing.T, hashAlg, dir string) {
	t.Helper()

	hash := ChooseHashAlg(hashAlg)
	require.NotNil(t, hash, "error hash not found")
	err := New(hash, config.Config{}, slog.New(slog.NewTextHandler(io.Discard, nil))).Calculate(dir)
	require.NoError(t, err, "error calculating hash")
}

func AppCheckHash(t *testing.T, hashAlg, dir string) {
	t.Helper()
	hash := ChooseHashAlg(hashAlg)
	require.NotNil(t, hash, "error hash not found")
	err := New(hash, config.Config{}, slog.New(slog.NewTextHandler(io.Discard, nil))).Check(dir)
	require.NoError(t, err, "error checking hash")
}

func CmdCalculateHastStdout(t *testing.T, appName string, data []byte) string {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd := exec.Command(appName)
	cmd.Stdin = bytes.NewReader(data)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	require.NoErrorf(t, err, "error calculating hash by application: %s", stderr.String())

	return stdout.String()
}

func CmdCalculateHash(t *testing.T, appName, hashAlg, dir string) {
	t.Helper()

	files, err := os.ReadDir(dir)
	require.NoErrorf(t, err, "error reading directory: %s", dir)
	filesForHash := make([]string, 0, len(files))
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		filesForHash = append(filesForHash, f.Name())
	}

	filePath := path.Join(dir, "checksum."+hashAlg)
	file, err := os.Create(filePath)
	require.NoErrorf(t, err, "error creating file: %s", filePath)
	defer func() {
		dErr := file.Sync()
		require.NoErrorf(t, dErr, "error syncing file: %s", filePath)
		dErr = file.Close()
		require.NoErrorf(t, dErr, "error closing file")
	}()
	stderr := &bytes.Buffer{}

	cmd := exec.Command(appName, filesForHash...)
	cmd.Dir = dir
	cmd.Stdout = file
	cmd.Stderr = stderr
	err = cmd.Run()
	require.NoErrorf(t, err, "error calculating hash by application: %s", stderr.String())
}

func CmdCheckHash(t *testing.T, appName, hashAlg, dir string, customFlags []string) {
	t.Helper()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	flags := []string{"--strict"}
	if len(customFlags) > 0 {
		flags = customFlags
	}
	flags = append(flags, "-c", "checksum."+hashAlg)

	cmd := exec.Command(appName, strings.Fields(strings.Join(flags, " "))...)
	cmd.Dir = dir
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	require.NoErrorf(t, err, "error check hash by application: %s", stderr.String())
}
