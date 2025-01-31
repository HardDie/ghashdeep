package crawler

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/HardDie/ghashdeep/internal/entities/checkfile"
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
func (c Crawler) calculateIterateFiles(checkPath string, filesForCalculate []string) error {
	onlyPath := filepath.Dir(checkPath)
	onlyDir := filepath.Base(checkPath)

	if len(filesForCalculate) == 0 {
		return nil
	}

	// Track time of checking current directory
	startedAt := time.Now()

	pr := utils.NewProcessing()

	check := checkfile.New(len(filesForCalculate))
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
				return fmt.Errorf("c.calculateFile: %w", err)
			}
			if Verbose {
				hashFinish = time.Now()
				c.logger.Debug(
					"stream hash calculation",
					slog.String(LogValueFile, fileName),
					slog.String(LogValueDuration, hashFinish.Sub(hashStart).String()),
				)
			}

			m.Lock()
			check.Add(fileName, fileHash)
			m.Unlock()
			return nil
		})
	}
	err := pr.Run()
	if err != nil {
		return fmt.Errorf("pr.Run: %w", err)
	}
	finishedAt := time.Now()

	if check.Len() > 0 {
		checkFilePath := path.Join(checkPath, c.checkFileName)
		err = check.SaveToFile(checkFilePath)
		if err != nil {
			return fmt.Errorf("check.SaveToFile: %w", err)
		}

		c.logger.Info(
			"Calculated!",
			slog.String(LogValuePath, onlyPath),
			slog.String(LogValueFolder, onlyDir),
			slog.String(LogValueDuration, finishedAt.Sub(startedAt).String()),
		)
	}

	return nil
}
