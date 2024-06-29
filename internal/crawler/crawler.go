package crawler

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

type HashMethod interface {
	Name() string
	Len() int
	Hash(file []byte) []byte
	Validate(file, hash []byte) bool
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
}

func New(hash HashMethod) *Crawler {
	return &Crawler{
		hash:          hash,
		checkFileName: "checksum." + hash.Name(),
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
