package crawler

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/HardDie/LibraryHashCheck/internal/utils"
)

func (c Crawler) Calculate(checkPath string) error {
	err := c.calculateIterate(checkPath)
	if err != nil {
		return err
	}
	return nil
}

func (c Crawler) calculateIterate(checkPath string) error {
	files, dirs, err := c.readFiles(checkPath)
	if err != nil {
		return err
	}

	info := make([]CheckFileInfo, 0, len(files))
	for _, file := range files {
		fileName := file.Name()
		if fileName == c.checkFileName {
			continue
		}
		fileData, err := utils.ReadAllFile(path.Join(checkPath, fileName))
		if err != nil {
			return fmt.Errorf("Crawler.calculateIterate(%s): %w", checkPath, err)
		}
		fileHash := c.hash.Hash(fileData)
		info = append(info, CheckFileInfo{
			Name:       fileName,
			Hash:       fileHash,
			HashString: hex.EncodeToString(fileHash),
		})
	}

	if len(info) > 0 {
		checkFilePath := path.Join(checkPath, c.checkFileName)
		err = c.writeCheckFile(checkFilePath, info)
		if err != nil {
			return fmt.Errorf("Crawler.calculateIterate(%s): %w", checkPath, err)
		}
	}

	for _, dir := range dirs {
		if dir.Name() == ".git" {
			// FIXME: remove
			continue
		}

		if err = c.calculateIterate(path.Join(checkPath, dir.Name())); err != nil {
			return nil
		}
	}

	return nil
}
func (c Crawler) writeCheckFile(checkFilePath string, info []CheckFileInfo) error {
	f, err := os.Create(checkFilePath)
	if err != nil {
		return fmt.Errorf("Crawler.writeCheckFile(%s) os.Create: %w", checkFilePath, err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			log.Printf("Crawler.writeCheckFile(%s) f.Close: %v", checkFilePath, e.Error())
		}
	}()

	for _, fileInfo := range info {
		_, err = f.WriteString(fileInfo.HashString + "  " + fileInfo.Name + "\n")
		if err != nil {
			return fmt.Errorf("Crawler.writeCheckFile(%s) f.WriteString: %w", checkFilePath, err)
		}
	}

	return nil
}
