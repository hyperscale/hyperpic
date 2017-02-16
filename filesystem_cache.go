package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"io/ioutil"

	"github.com/euskadi31/image-service/memfs"
	"github.com/rs/xlog"
)

type FileSystemCacheProvider struct {
	path string
}

func NewFileSystemCacheProvider(config *CacheFSConfiguration) *FileSystemCacheProvider {
	return &FileSystemCacheProvider{
		path: config.Path,
	}
}

func (p FileSystemCacheProvider) Get(file string, options *ImageOptions) (*Resource, bool) {
	if containsDotDot(file) {
		xlog.Error("Invalid URL path")

		return nil, false
	}

	path := p.path + "/" + strings.TrimPrefix(file, "/") + "/" + options.Hash()

	f, err := os.Open(path)
	if os.IsNotExist(err) {
		xlog.Debugf("File %s is not found in cache", file)

		return nil, false
	}

	if err != nil {
		xlog.Error("Cannot open file")

		return nil, false
	}
	defer f.Close()

	d, err := f.Stat()
	if err != nil {
		xlog.Error("Cannot read file info")

		return nil, false
	}

	if d.IsDir() {
		xlog.Error("Is not a file")

		return nil, false
	}

	body, err := ioutil.ReadAll(f)
	if err != nil {
		xlog.Error("Cannot read file")

		return nil, false
	}

	_, name := filepath.Split(path)

	return &Resource{
		Path:       file,
		Name:       name,
		Options:    options,
		Body:       body,
		ModifiedAt: d.ModTime(),
	}, true
}

func (p FileSystemCacheProvider) Set(resource *Resource) error {
	path := p.path + "/" + strings.TrimPrefix(resource.Path, "/")
	filename := path + "/" + resource.Options.Hash()

	if err := os.MkdirAll(path, 0777); err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	n, err := io.Copy(file, memfs.NewBuffer(&resource.Body))
	if err != nil {
		return err
	}

	xlog.Debugf("Write cache size: %d", n)

	return nil
}
