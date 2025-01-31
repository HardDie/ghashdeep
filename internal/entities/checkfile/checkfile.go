package checkfile

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/HardDie/ghashdeep/internal/utils"
)

type CheckFile struct {
	Files map[string]Object
}

func New(size int) *CheckFile {
	return &CheckFile{
		Files: make(map[string]Object, size),
	}
}

func NewFromFile(path string, hashLen int) (*CheckFile, error) {
	data, err := utils.ReadAllFile(path)
	if err != nil {
		return nil, fmt.Errorf("utils.ReadAllFile(%s): %w", path, err)
	}

	lines := strings.Split(string(data), "\n")
	res := make(map[string]Object, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}
		// two spaces and at least one name symbol
		if len(line) < hashLen+3 {
			return nil, fmt.Errorf("(%q) invalid length", line)
		}

		obj, err := NewObjectFromString(line, hashLen)
		if err != nil {
			return nil, fmt.Errorf("NewObjectFromString: %w", err)
		}
		res[obj.Name] = obj
	}

	return &CheckFile{
		Files: res,
	}, nil
}

func (c CheckFile) Len() int {
	return len(c.Files)
}

func (c CheckFile) IsFileExist(name string) (Object, bool) {
	obj, ok := c.Files[name]
	return obj, ok
}

func (c CheckFile) SaveToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("os.Create(%s): %w", path, err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			log.Printf("file: %s f.Close: %v", path, e.Error())
		}
	}()

	info := make([]Object, 0, len(c.Files))
	for _, obj := range c.Files {
		info = append(info, obj)
	}
	sort.SliceStable(info, func(i, j int) bool {
		return info[i].Name < info[j].Name
	})

	for _, fileInfo := range info {
		_, err = f.WriteString(fileInfo.String() + "\n")
		if err != nil {
			return fmt.Errorf("file: %s f.WriteString: %w", path, err)
		}
	}

	return nil
}

func (c *CheckFile) Delete(name string) {
	delete(c.Files, name)
}

func (c *CheckFile) Add(name string, hash []byte) {
	c.Files[name] = Object{
		Name: name,
		Hash: hash,
	}
}
