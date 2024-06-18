package crawler

import (
	"encoding/hex"
	"fmt"
	"log"
	"log/slog"
	"path"
	"strings"
	"time"

	"github.com/HardDie/LibraryHashCheck/internal/logger"
	"github.com/HardDie/LibraryHashCheck/internal/utils"
)

func (c Crawler) Check(checkPath string) error {
	err := c.checkIterate(checkPath)
	if err != nil {
		return err
	}
	return nil
}

func (c Crawler) checkIterate(checkPath string) error {
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

	if len(filesForCheck) > 0 {
		if !isCheckFileExist {
			logger.Error(
				"no checksum",
				slog.String(logger.LogValueStatus, "BAD"),
				slog.String(logger.LogValuePath, checkPath),
			)
		} else {
			startedAt := time.Now()
			info, err := c.readCheckFile(path.Join(checkPath, c.checkFileName))
			if err != nil {
				log.Println("err", err.Error())
				return err
			}

			badFiles := make([]string, 0, len(filesForCheck))
			notFound := make([]string, 0, len(filesForCheck))
			for _, fileName := range filesForCheck {
				fileData, err := utils.ReadAllFile(path.Join(checkPath, fileName))
				if err != nil {
					return fmt.Errorf("Crawler.iterate(%s): %w", checkPath, err)
				}
				fileInfo, ok := info[fileName]
				if !ok {
					notFound = append(notFound, fileName)
					continue
				}
				// Exclude duplication
				delete(info, fileName)
				if !c.hash.Validate(fileData, fileInfo.Hash) {
					badFiles = append(badFiles, fileName)
				}
			}
			finishedAt := time.Now()

			if len(badFiles) > 0 ||
				len(notFound) > 0 ||
				len(info) > 0 {
				logger.Error(
					"folder have errors",
					slog.String(logger.LogValueStatus, "BAD"),
					slog.String(logger.LogValuePath, checkPath),
					slog.String(logger.LogValueStartedAt, startedAt.String()),
					slog.String(logger.LogValueFinishedAt, finishedAt.String()),
					slog.String(logger.LogValueDuration, finishedAt.Sub(startedAt).String()),
				)
				for _, badFile := range notFound {
					logger.Error(
						"no checksum",
						slog.String(logger.LogValueStatus, "BAD"),
						slog.String(logger.LogValueFile, badFile),
					)
				}
				for _, badFile := range badFiles {
					logger.Error(
						"bad checksum",
						slog.String(logger.LogValueStatus, "BAD"),
						slog.String(logger.LogValueFile, badFile),
					)
				}
				for _, fileInfo := range info {
					logger.Error(
						"not found",
						slog.String(logger.LogValueStatus, "BAD"),
						slog.String(logger.LogValueFile, fileInfo.Name),
					)
				}
			} else {
				logger.Info(
					"Success",
					slog.String(logger.LogValueStatus, "GOOD"),
					slog.String(logger.LogValuePath, checkPath),
					slog.String(logger.LogValueStartedAt, startedAt.String()),
					slog.String(logger.LogValueFinishedAt, finishedAt.String()),
					slog.String(logger.LogValueDuration, finishedAt.Sub(startedAt).String()),
				)
			}
		}
	}

	for _, dir := range dirs {
		if dir.Name() == ".git" {
			// FIXME: remove
			continue
		}

		if err = c.checkIterate(path.Join(checkPath, dir.Name())); err != nil {
			return nil
		}
	}

	return nil
}

func (c Crawler) splitChecksumFileLine(line string) (CheckFileInfo, error) {
	if len(line) < c.hash.Len()+3 /* two spaces and at least one name symbol */ {
		return CheckFileInfo{}, fmt.Errorf("invalid length")
	}

	info := CheckFileInfo{
		HashString: line[:c.hash.Len()],
		Name:       line[c.hash.Len()+2:],
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
	//res := make([]CheckFileInfo, 0, len(lines))
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
