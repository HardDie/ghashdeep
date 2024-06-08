package utils

import (
	"fmt"
	"io"
	"log"
	"os"
)

func ReadAllFile(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("utils.ReadAllFile(%s) os.Open: %w", filePath, err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			log.Printf("utils.ReadAllFile(%s) f.Close: %v", filePath, e.Error())
		}
	}()
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("utils.ReadAllFile(%s) io.ReadAll: %w", filePath, err)
	}
	return data, nil
}
