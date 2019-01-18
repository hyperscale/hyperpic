// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package filesystem

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hyperscale/hyperpic/pkg/hyperpic/fsutil"
	"github.com/hyperscale/hyperpic/pkg/hyperpic/image"
	"github.com/hyperscale/hyperpic/pkg/hyperpic/memfs"
	"github.com/rs/zerolog/log"
)

// CacheProvider struct
type CacheProvider struct {
	config *CacheConfiguration
}

// NewCacheProvider constructor of FS Cache provider
func NewCacheProvider(cfg *CacheConfiguration) *CacheProvider {
	p := &CacheProvider{
		config: cfg,
	}

	p.Run()

	return p
}

func (p CacheProvider) removeOldCacheFile(path string, f os.FileInfo, err error) error {
	if f == nil {
		return nil
	}

	now := time.Now()

	if !f.IsDir() && now.After(f.ModTime().Add(p.config.LifeTime)) {
		log.Debug().Msgf("Remove file %s", path)
		if err := os.Remove(path); err != nil {
			return err
		}
	}

	return nil
}

// Run cleanner
func (p CacheProvider) Run() {
	log.Debug().Msg("Cleanner running")

	ticker := time.NewTicker(p.config.CleanInterval)

	go func() {
		for range ticker.C {
			log.Debug().Msg("Start cleaning")

			if err := filepath.Walk(p.config.Path, p.removeOldCacheFile); err != nil {
				log.Error().Err(err).Msgf("Walk on %s", p.config.Path)
			}

			log.Debug().Msg("End cleaning")
		}
	}()
}

// Del all cache files for source file
func (p CacheProvider) Del(resource *image.Resource) error {
	if fsutil.ContainsDotDot(resource.Path) {
		return errors.New("Invalid URL path")
	}

	path := p.config.Path + "/" + strings.TrimPrefix(resource.Path, "/")

	return os.RemoveAll(path)
}

// Get cached file
func (p CacheProvider) Get(resource *image.Resource) (*image.Resource, error) {
	if fsutil.ContainsDotDot(resource.Path) {
		return nil, ErrInvalidPath
	}

	path := p.config.Path + "/" + strings.TrimPrefix(resource.Path, "/") + "/" + resource.Options.Hash()

	f, err := os.Open(path)
	if os.IsNotExist(err) {
		log.Debug().Msgf("File %s is not found in cache", resource.Path)

		return nil, ErrCacheNotExist
	}

	if err != nil {
		return nil, err
	}
	defer f.Close()

	d, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if d.IsDir() {
		log.Debug().Msgf("%s is not a file", resource.Path)

		return nil, ErrNotFile
	}

	body, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	_, name := filepath.Split(resource.Path)

	return &image.Resource{
		Path:       resource.Path,
		Name:       name,
		Options:    resource.Options,
		Body:       body,
		Size:       len(body),
		ModifiedAt: d.ModTime(),
	}, nil
}

// Set file to cache
func (p CacheProvider) Set(resource *image.Resource) error {
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

	log.Debug().Msgf("Write cache file size: %d", n)

	return nil
}
