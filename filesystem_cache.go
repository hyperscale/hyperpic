// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"io/ioutil"

	"github.com/euskadi31/image-service/memfs"
	"github.com/rs/xlog"
)

type FileSystemCacheProvider struct {
	config *CacheFSConfiguration
}

// NewFileSystemCacheProvider constructor of FS Cache provider
func NewFileSystemCacheProvider(config *CacheFSConfiguration) *FileSystemCacheProvider {
	p := &FileSystemCacheProvider{
		config: config,
	}

	p.Run()

	return p
}

func (p *FileSystemCacheProvider) removeOldCacheFile(path string, f os.FileInfo, err error) error {
	now := time.Now()

	if f.IsDir() == false && now.After(f.ModTime().Add(p.config.LifeTime)) {
		xlog.Debugf("Remove file %s", path)
		if err := os.Remove(path); err != nil {
			xlog.Error(err)
		}
	}

	return nil
}

func (p *FileSystemCacheProvider) Run() {
	xlog.Debug("Cleanner running")
	ticker := time.NewTicker(p.config.CleanInterval)
	// defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				xlog.Debug("Start cleaning")
				err := filepath.Walk(p.config.Path, p.removeOldCacheFile)
				if err != nil {
					xlog.Error(err)
				}
				xlog.Debug("End cleaning")
			}
		}
	}()
}

// Del all cache files for source file
func (p FileSystemCacheProvider) Del(resource *Resource) bool {
	if containsDotDot(resource.Path) {
		xlog.Error("Invalid URL path")

		return false
	}
	path := p.config.Path + "/" + strings.TrimPrefix(resource.Path, "/")

	if err := os.RemoveAll(path); err != nil {
		return false
	}

	return true
}

// Get cached file
func (p FileSystemCacheProvider) Get(resource *Resource) (*Resource, bool) {
	if containsDotDot(resource.Path) {
		xlog.Error("Invalid URL path")

		return nil, false
	}

	path := p.config.Path + "/" + strings.TrimPrefix(resource.Path, "/") + "/" + resource.Options.Hash()

	f, err := os.Open(path)
	if os.IsNotExist(err) {
		xlog.Debugf("File %s is not found in cache", resource.Path)

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

	_, name := filepath.Split(resource.Path)

	return &Resource{
		Path:       resource.Path,
		Name:       name,
		Options:    resource.Options,
		Body:       body,
		ModifiedAt: d.ModTime(),
	}, true
}

// Set file to cache
func (p FileSystemCacheProvider) Set(resource *Resource) error {
	path := p.config.Path + "/" + strings.TrimPrefix(resource.Path, "/")
	filename := path + "/" + resource.Options.Hash()

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
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

	xlog.Debugf("Write cache file size: %d", n)

	return nil
}
