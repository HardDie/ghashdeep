package crawler

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"github.com/HardDie/ghashdeep/internal/validators/blake3"
	"github.com/HardDie/ghashdeep/internal/validators/md5"
	"github.com/HardDie/ghashdeep/internal/validators/sha1"
	"github.com/HardDie/ghashdeep/internal/validators/sha224"
	"github.com/HardDie/ghashdeep/internal/validators/sha256"
	"github.com/HardDie/ghashdeep/internal/validators/sha384"
	"github.com/HardDie/ghashdeep/internal/validators/sha512"
	"github.com/HardDie/ghashdeep/internal/validators/xxhash"
)

var (
	ErrChecksumNotFound = errors.New("checksum file not found")
	ErrHaveInvalidFiles = errors.New("have invalid files")
)

type HashMethod interface {
	Name() string
	CalculateStream(s io.Reader) ([]byte, error)
}

type CheckFileInfo struct {
	Name       string
	HashString string
	Hash       []byte
}

type Crawler struct {
	hash          HashMethod
	checkFileName string
	hashLen       int
}

func New(hash HashMethod) *Crawler {
	return &Crawler{
		hash:          hash,
		checkFileName: "checksum." + hash.Name(),
		hashLen:       calcHashLen(hash),
	}
}

func (c Crawler) ValidateStream(s io.Reader, hash []byte) (bool, error) {
	fileHash, err := c.hash.CalculateStream(s)
	if err != nil {
		return false, fmt.Errorf("CalculateStream: %w", err)
	}
	return bytes.Equal(fileHash, hash), nil
}

func (c Crawler) readFiles(checkPath string) ([]os.FileInfo, []os.FileInfo, error) {
	f, err := os.Open(checkPath)
	if err != nil {
		return nil, nil, fmt.Errorf("Crawler.readFiles(%s) os.Open: %w", checkPath, err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			log.Printf("Crawler.readFiles(%s) f.Close: %v", checkPath, e.Error())
		}
	}()

	st, err := f.Stat()
	if err != nil {
		return nil, nil, fmt.Errorf("Crawler.readFiles(%s) f.Stat: %w", checkPath, err)
	}
	if !st.IsDir() {
		return nil, nil, fmt.Errorf("Crawler.readFiles(%s) passed path is not directory", checkPath)
	}

	files, err := f.Readdir(0)
	if err != nil {
		return nil, nil, fmt.Errorf("Crawler.readFiles(%s) f.Readdir: %w", checkPath, err)
	}

	var resFiles []os.FileInfo
	var resDirs []os.FileInfo

	for _, file := range files {
		if file.IsDir() {
			resDirs = append(resDirs, file)
			continue
		}
		resFiles = append(resFiles, file)
	}

	sort.SliceStable(resFiles, func(i, j int) bool {
		return resFiles[i].Name() < resFiles[j].Name()
	})
	sort.SliceStable(resDirs, func(i, j int) bool {
		return resDirs[i].Name() < resDirs[j].Name()
	})

	return resFiles, resDirs, nil
}

func calcHashLen(
	hash interface {
		CalculateStream(s io.Reader) ([]byte, error)
	},
) int {
	fileHash, err := hash.CalculateStream(bytes.NewReader([]byte{1}))
	if err != nil {
		panic(err)
	}
	return len(hex.EncodeToString(fileHash))
}

func ChooseHashAlg(alg string) HashMethod {
	if alg == "" {
		alg = "md5"
	}
	switch alg {
	case "md5":
		return md5.New()
	case "sha1":
		return sha1.New()
	case "sha224":
		return sha224.New()
	case "sha256":
		return sha256.New()
	case "sha384":
		return sha384.New()
	case "sha512":
		return sha512.New()
	case "xxhash":
		return xxhash.New()
	case "blake3":
		return blake3.New()
	}
	return nil
}
