package main

import (
	"fmt"
	"os"
	"strings"
)

type FileSystemSourceProvider struct {
	path string
}

func NewFileSystemSourceProvider(path string) *FileSystemSourceProvider {
	return &FileSystemSourceProvider{
		path: path,
	}
}

func (p FileSystemSourceProvider) Get(file string) (*Resource, error) {
	path := p.path + "/" + strings.TrimPrefix(file, "/")

	finfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("File %s not exists", file)
	}

	if finfo.IsDir() {
		return nil, fmt.Errorf("%s is not a file", file)
	}

	return &Resource{
		File:       file,
		SourcePath: path,
	}, nil
}
