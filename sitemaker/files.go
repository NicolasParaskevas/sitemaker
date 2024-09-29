package sitemaker

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Filetype int

const (
	DATA Filetype = iota
	ASSET
)

func loadSourceFiles(source string) (map[string]string, error) {
	d, err := load(source, DATA)

	if err != nil {
		return nil, err
	}

	data := make(map[string]string)

	for k, v := range d {
		data[k] = string(v)
	}

	return data, err
}

func loadAssetFiles(source string) (map[string][]byte, error) {
	assetDir := filepath.Join(source, "assets")
	return load(assetDir, ASSET)
}

func load(source string, ft Filetype) (map[string][]byte, error) {
	result := make(map[string][]byte)

	err := filepath.Walk(source, func(fpath string, info os.FileInfo, err error) error {

		if !info.IsDir() {
			b, err := os.ReadFile(fpath)

			if err != nil {
				return err
			}

			if !strings.HasPrefix(fpath, ".") && !strings.Contains(fpath, "/.") {
				if ft == DATA && strings.HasSuffix(fpath, ".xml") {
					result[fpath] = b
				}

				if ft == ASSET && !strings.HasSuffix(fpath, ".xml") {
					result[fpath] = b
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, err
}

func makeDirs(filename string, mode os.FileMode) error {
	baseDir := path.Dir(filename)
	info, err := os.Stat(baseDir)
	if err == nil && info.IsDir() {
		return nil
	}
	return os.MkdirAll(baseDir, mode)
}
