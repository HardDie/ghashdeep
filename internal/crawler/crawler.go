package crawler

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"github.com/HardDie/ghashdeep/internal/validators"
)

type HashMethod interface {
	Name() string
	Hash(file []byte) []byte
	CalculateStream(s io.Reader) ([]byte, error)
	ValidateStream(s io.Reader, hash []byte) (bool, error)
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
		Hash(file []byte) []byte
	},
) int {
	return len(hex.EncodeToString(hash.Hash([]byte{1})))
}

func ChooseHashAlg(alg string) HashMethod {
	if alg == "" {
		alg = "md5"
	}
	switch alg {
	case "md5":
		return validators.NewMd5()
	case "sha1":
		return validators.NewSha1()
	case "sha224":
		return validators.NewSha224()
	case "sha256":
		return validators.NewSha256()
	case "sha384":
		return validators.NewSha384()
	case "sha512":
		return validators.NewSha512()
	case "xxhash":
		return validators.NewXxhash()
	case "blake3":
		return validators.NewBlake3()
	}
	return nil
}
