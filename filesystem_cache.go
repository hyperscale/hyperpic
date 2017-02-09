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

func (p FileSystemCacheProvider) Get(file string, options ImageOptions) (*Resource, bool) {
	path := p.path + "/" + strings.TrimPrefix(file, "/") + "/" + options.Hash()

	finfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, false
	}

	if finfo.IsDir() {
		return nil, false
	}

	return nil, true
}

func (p FileSystemCacheProvider) Set(resource *Resource) error {

	return nil
}
