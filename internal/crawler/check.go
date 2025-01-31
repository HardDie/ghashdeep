package crawler

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/HardDie/ghashdeep/internal/utils"
)

const (
	Verbose = false
)

func (c Crawler) Check(checkPath string) error {
	err := c.checkIterate(checkPath)
	if err != nil {
		return err
	}
	return nil
}

func (c Crawler) checkIterate(checkPath string) error {
	var storeError error

	files, dirs, err := c.readFiles(checkPath)
	if err != nil {
		return err
	}

	filesForCheck := make([]string, 0, len(files))
	isCheckFileExist := false
	for _, file := range files {
		fileName := file.Name()
		if fileName == c.checkFileName {
			isCheckFileExist = true
			continue
		}
		filesForCheck = append(filesForCheck, fileName)
	}

	// Check files
	err = c.checkIterateFiles(checkPath, isCheckFileExist, filesForCheck)
	if err != nil {
		switch {
		case errors.Is(err, ErrChecksumNotFound):
			storeError = err
		case errors.Is(err, ErrHaveInvalidFiles):
			storeError = err
		default:
			return fmt.Errorf("Crawler.checkIterate(%s): %w", checkPath, err)
		}
	}

	// Recursive check other directories
	for _, dir := range dirs {
		if dir.Name() == ".git" {
			// FIXME: remove
			continue
		}

		if err = c.checkIterate(path.Join(checkPath, dir.Name())); err != nil {
			switch {
			case errors.Is(err, ErrChecksumNotFound):
				storeError = err
			case errors.Is(err, ErrHaveInvalidFiles):
				storeError = err
			default:
				return nil
			}
		}
	}

	if storeError != nil {
		return storeError
	}
	return nil
}
func (c Crawler) validateFile(checkFilePath string, hash []byte) (bool, error) {
	f, err := os.Open(checkFilePath)
	if err != nil {
		return false, fmt.Errorf("Crawler.validateFile(%s) os.Open: %w", checkFilePath, err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			log.Printf("Crawler.validateFile(%s) f.Close: %v", checkFilePath, e.Error())
		}
	}()

	isValid, err := c.ValidateStream(f, hash)
	if err != nil {
		return false, fmt.Errorf("Crawler.validateFile(%s) ValidateStream: %w", checkFilePath, err)
	}

	return isValid, nil
}
func (c Crawler) checkIterateFiles(checkPath string, isCheckFileExist bool, filesForCheck []string) error {
	onlyPath := filepath.Dir(checkPath)
	onlyDir := filepath.Base(checkPath)

	if len(filesForCheck) == 0 && !isCheckFileExist {
		return nil
	}

	if !isCheckFileExist {
		c.logger.Error(
			"no checksum",
			slog.String(LogValueStatus, "BAD"),
			slog.String(LogValuePath, checkPath),
		)
		return ErrChecksumNotFound
	}

	// Track time of checking current directory
	startedAt := time.Now()

	// Parse file with validation info
	info, err := c.readCheckFile(path.Join(checkPath, c.checkFileName))
	if err != nil {
		log.Println("err", err.Error())
		return err
	}

	// Prepare lists for files with invalid checksums and not exist files
	var badFiles, notFound []string

	pr := utils.NewProcessing()

	// Iterate through all files and compare checksum
	for _, fileName := range filesForCheck {
		// Check if file exist in checksum file
		fileInfo, ok := info[fileName]
		if !ok {
			notFound = append(notFound, fileName)
			continue
		}
		// Exclude duplication
		delete(info, fileName)

		// Build full path to required file
		fullFilePath := path.Join(checkPath, fileName)

		pr.Push(fileName, func(m *sync.Mutex) error {
			// Track time calculation of hash sum for selected file
			var hashStart, hashFinish time.Time
			if Verbose {
				hashStart = time.Now()
			}

			isValid, err := c.validateFile(fullFilePath, fileInfo.Hash)
			if err != nil {
				return fmt.Errorf("Crawler.checkIterateFiles: %w", err)
			}
			if Verbose {
				hashFinish = time.Now()
				c.logger.Debug(
					"stream hash calculation",
					slog.String(LogValueFile, fileName),
					slog.String(LogValueDuration, hashFinish.Sub(hashStart).String()),
				)
			}
			if !isValid {
				m.Lock()
				badFiles = append(badFiles, fileName)
				m.Unlock()
			}
			return nil
		})
	}
	err = pr.Run()
	if err != nil {
		return fmt.Errorf("Crawler.checkIterateFiles() pr.Run: %w", err)
	}
	finishedAt := time.Now()

	// Print all found errors for current directory
	if len(badFiles) > 0 ||
		len(notFound) > 0 ||
		len(info) > 0 {
		c.logger.Error(
			"folder have errors",
			slog.String(LogValueStatus, "BAD"),
			slog.String(LogValuePath, onlyPath),
			slog.String(LogValueFolder, onlyDir),
			// slog.String(LogValueStartedAt, startedAt.String()),
			// slog.String(LogValueFinishedAt, finishedAt.String()),
			slog.String(LogValueDuration, finishedAt.Sub(startedAt).String()),
		)
		for _, badFile := range notFound {
			c.logger.Error(
				"no checksum",
				slog.String(LogValueStatus, "BAD"),
				slog.String(LogValueFile, badFile),
			)
		}
		for _, badFile := range badFiles {
			c.logger.Error(
				"bad checksum",
				slog.String(LogValueStatus, "BAD"),
				slog.String(LogValueFile, badFile),
			)
		}
		for _, fileInfo := range info {
			c.logger.Error(
				"not found",
				slog.String(LogValueStatus, "BAD"),
				slog.String(LogValueFile, fileInfo.Name),
			)
		}
		return ErrHaveInvalidFiles
	}

	c.logger.Info(
		"Success",
		slog.String(LogValueStatus, "GOOD"),
		slog.String(LogValuePath, onlyPath),
		slog.String(LogValueFolder, onlyDir),
		// slog.String(LogValueStartedAt, startedAt.String()),
		// slog.String(LogValueFinishedAt, finishedAt.String()),
		slog.String(LogValueDuration, finishedAt.Sub(startedAt).String()),
	)

	return nil
}

func (c Crawler) splitChecksumFileLine(line string) (CheckFileInfo, error) {
	if len(line) < c.hashLen+3 /* two spaces and at least one name symbol */ {
		return CheckFileInfo{}, fmt.Errorf("invalid length")
	}

	info := CheckFileInfo{
		HashString: line[:c.hashLen],
		Name:       line[c.hashLen+2:],
	}
	hash, err := hex.DecodeString(info.HashString)
	if err != nil {
		return CheckFileInfo{}, fmt.Errorf("hex.DecodeString: %w", err)
	}
	info.Hash = hash
	return info, nil
}
func (c Crawler) readCheckFile(checkFilePath string) (map[string]CheckFileInfo, error) {
	data, err := utils.ReadAllFile(checkFilePath)
	if err != nil {
		return nil, fmt.Errorf("Crawler.readCheckFile(%s): %w", checkFilePath, err)
	}

	lines := strings.Split(string(data), "\n")
	res := make(map[string]CheckFileInfo, len(lines))
	// res := make([]CheckFileInfo, 0, len(lines))
	for i, line := range lines {
		if len(line) == 0 {
			continue
		}
		info, err := c.splitChecksumFileLine(line)
		if err != nil {
			return nil, fmt.Errorf("Crawler.readCheckFile(%s) line %d %w", checkFilePath, i, err)
		}
		res[info.Name] = info
	}

	return res, nil
}
