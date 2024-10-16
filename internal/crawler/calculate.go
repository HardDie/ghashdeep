package crawler

import (
	"encoding/hex"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/HardDie/ghashdeep/internal/logger"
	"github.com/HardDie/ghashdeep/internal/utils"
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

	filesForCalculate := make([]string, 0, len(files))
	for _, file := range files {
		fileName := file.Name()
		if fileName == c.checkFileName {
			continue
		}
		filesForCalculate = append(filesForCalculate, fileName)
	}

	// Calculate hash for files
	err = c.calculateIterateFiles(checkPath, filesForCalculate)
	if err != nil {
		return fmt.Errorf("Crawler.calculateIterate(%s): %w", checkPath, err)
	}

	// Recursive check other directories
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
func (c Crawler) calculateFile(calculateFilePath string) ([]byte, error) {
	f, err := os.Open(calculateFilePath)
	if err != nil {
		return nil, fmt.Errorf("Crawler.calculateFile(%s) os.Open: %w", calculateFilePath, err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			log.Printf("Crawler.calculateFile(%s) f.Close: %v", calculateFilePath, e.Error())
		}
	}()

	hash, err := c.hash.CalculateStream(f)
	if err != nil {
		return nil, fmt.Errorf("Crawler.calculateFile(%s) CalculateStream: %w", calculateFilePath, err)
	}

	return hash, nil
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
func (c Crawler) calculateIterateFiles(checkPath string, filesForCalculate []string) error {
	onlyPath := filepath.Dir(checkPath)
	onlyDir := filepath.Base(checkPath)

	if len(filesForCalculate) == 0 {
		return nil
	}

	// Track time of checking current directory
	startedAt := time.Now()

	pr := utils.NewProcessing()

	info := make([]CheckFileInfo, 0, len(filesForCalculate))
	for _, fileName := range filesForCalculate {
		// Build full path to required file
		fullFilePath := path.Join(checkPath, fileName)

		pr.Push(fileName, func(m *sync.Mutex) error {
			// Track time calculation of hash sum for selected file
			var hashStart, hashFinish time.Time
			if Verbose {
				hashStart = time.Now()
			}

			fileHash, err := c.calculateFile(fullFilePath)
			if err != nil {
				return fmt.Errorf("Crawler.calculateIterateFiles: %w", err)
			}
			if Verbose {
				hashFinish = time.Now()
				logger.Debug(
					"stream hash calculation",
					slog.String(logger.LogValueFile, fileName),
					slog.String(logger.LogValueDuration, hashFinish.Sub(hashStart).String()),
				)
			}

			m.Lock()
			info = append(info, CheckFileInfo{
				Name:       fileName,
				Hash:       fileHash,
				HashString: hex.EncodeToString(fileHash),
			})
			m.Unlock()
			return nil
		})
	}
	err := pr.Run()
	if err != nil {
		return fmt.Errorf("Crawler.calculateIterateFiles() pr.Run: %w", err)
	}
	finishedAt := time.Now()

	if len(info) > 0 {
		sort.SliceStable(info, func(i, j int) bool {
			return info[i].Name < info[j].Name
		})

		checkFilePath := path.Join(checkPath, c.checkFileName)
		err = c.writeCheckFile(checkFilePath, info)
		if err != nil {
			return fmt.Errorf("Crawler.calculateIterateFiles(%s): %w", checkPath, err)
		}
		logger.Info(
			"Calculated!",
			slog.String(logger.LogValuePath, onlyPath),
			slog.String(logger.LogValueFolder, onlyDir),
			slog.String(logger.LogValueDuration, finishedAt.Sub(startedAt).String()),
		)
	}

	return nil
}
