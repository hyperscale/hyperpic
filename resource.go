package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// Resource struct
type Resource struct {
	File       string
	SourcePath string
	CachePath  string
	MimeType   string
	Body       []byte
}

// NewResource constructor
func NewResource(config *ImageConfiguration, file string) (*Resource, error) {
	sourcePath := config.SourcePathWithFile(file)

	finfo, err := os.Stat(sourcePath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("File %s not exists", file)
	}

	if finfo.IsDir() {
		return nil, fmt.Errorf("%s is not a file", file)
	}

	return &Resource{
		File:       file,
		SourcePath: sourcePath,
		CachePath:  config.CachePathWithFile(file),
	}, nil
}

func (r *Resource) Read() ([]byte, error) {
	return ioutil.ReadFile(r.SourcePath)
}
