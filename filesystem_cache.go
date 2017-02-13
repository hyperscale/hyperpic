package main

import (
	"os"
	"strings"
)

type FileSystemCacheProvider struct {
	path string
}

func NewFileSystemCacheProvider(path string) *FileSystemCacheProvider {
	return &FileSystemCacheProvider{
		path: path,
	}
}

func (p FileSystemCacheProvider) Get(file string, options *ImageOptions) (*Resource, bool) {
	path := p.path + "/" + strings.TrimPrefix(file, "/") + "/" + options.Hash()

	finfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, false
	}

	if finfo.IsDir() {
		return nil, false
	}

	return &Resource{
		CachePath: path,
	}, true
}

func (p FileSystemCacheProvider) Set(resource *Resource, options *ImageOptions) error {
	path := p.path + "/" + strings.TrimPrefix(resource.File, "/")
	filename := path + "/" + options.Hash()

	if err := os.MkdirAll(path, 0777); err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(resource.Body); err != nil {
		return err
	}

	return nil
}
