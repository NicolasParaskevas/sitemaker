package sitemaker

import (
	"os"
	"path/filepath"
	"strings"
)

func loadSourceFiles(source string) (map[string]string, map[string][]byte, error) {
	data := make(map[string]string)
	assets := make(map[string][]byte)

	err := filepath.Walk(source, func(fpath string, info os.FileInfo, err error) error {

		if !info.IsDir() {
			b, err := os.ReadFile(fpath)

			if err != nil {
				return err
			}

			// avoid files with . prefix
			if !strings.HasPrefix(fpath, ".") &&
				!strings.Contains(fpath, "/.") {
				if strings.HasSuffix(fpath, ".txt") {
					data[fpath] = string(b)
				} else {
					assets[fpath] = b
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return data, assets, err
}
