// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/xlog"
)

type FileSystemSourceProvider struct {
	path string
}

func NewFileSystemSourceProvider(config *SourceFSConfiguration) *FileSystemSourceProvider {
	return &FileSystemSourceProvider{
		path: config.Path,
	}
}

func (p FileSystemSourceProvider) Set(resource *Resource) error {
	path := p.path + "/" + strings.TrimPrefix(resource.Path, "/")

	dir, _ := filepath.Split(path)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	n, err := file.Write(resource.Body)
	if err != nil {
		return err
	}

	xlog.Debugf("Write source file size: %d", n)

	return nil
}

func (p FileSystemSourceProvider) Get(resource *Resource) (*Resource, error) {
	if containsDotDot(resource.Path) {
		// Too many programs use r.URL.Path to construct the argument to
		// serveFile. Reject the request under the assumption that happened
		// here and ".." may not be wanted.
		// Note that name might not contain "..", for example if code (still
		// incorrectly) used filepath.Join(myDir, r.URL.Path).

		return nil, fmt.Errorf("Invalid URL path")
	}

	path := p.path + "/" + strings.TrimPrefix(resource.Path, "/")

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	d, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if d.IsDir() {
		return nil, fmt.Errorf("%s is not a file", resource.Path)
	}

	body, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	_, name := filepath.Split(resource.Path)

	return &Resource{
		Path:       resource.Path,
		Options:    resource.Options,
		Name:       name,
		Body:       body,
		ModifiedAt: d.ModTime(),
	}, nil
}

// Del source files
func (p FileSystemSourceProvider) Del(resource *Resource) bool {
	if containsDotDot(resource.Path) {
		xlog.Error("Invalid URL path")

		return false
	}
	path := p.path + "/" + strings.TrimPrefix(resource.Path, "/")

	if err := os.Remove(path); err != nil {
		return false
	}

	return true
}
