package main

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
)

//go:embed *
var fsys embed.FS

func main() {
	result, err := LoadConfigs(fsys, ".")
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}

func LoadConfigs(fsys fs.FS, root string) (map[string][]byte, error) {
	result := make(map[string][]byte)
	fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if filepath.Ext(path) == ".conf" {
			byte, err := fs.ReadFile(fsys, path)
			if err != nil {
				return err
			}
			result[path] = byte
		}

		return nil
	})

	return result, nil
}
