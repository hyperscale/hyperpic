// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package memory

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hyperscale/hyperpic/pkg/hyperpic/fsutil"
	"github.com/hyperscale/hyperpic/pkg/hyperpic/image"
	"github.com/rs/zerolog/log"
)

// CacheProvider struct
type CacheProvider struct {
	config    *CacheConfiguration
	size      uint64
	mtx       sync.RWMutex
	container map[string]map[string]*image.Resource
}

// NewCacheProvider constructor of FS Cache provider
func NewCacheProvider(cfg *CacheConfiguration) *CacheProvider {
	p := &CacheProvider{
		config:    cfg,
		container: make(map[string]map[string]*image.Resource),
	}

	p.Run()

	return p
}

func (p *CacheProvider) removeOldCache(path string, key string, resource *image.Resource) {
	now := time.Now()

	if now.After(resource.ModifiedAt.Add(p.config.LifeTime)) {
		log.Debug().Msgf("Remove file %s", path)

		p.mtx.Lock()
		delete(p.container[path], key)
		p.mtx.Unlock()

		atomic.AddUint64(&p.size, ^uint64(resource.Size))
	}
}

// Run cleanner
func (p *CacheProvider) Run() {
	log.Debug().Msg("Cleanner running")

	ticker := time.NewTicker(p.config.CleanInterval)

	go func() {
		for range ticker.C {
			log.Debug().Msg("Start cleaning")

			for path, container := range p.container {
				for key, resource := range container {
					p.removeOldCache(path, key, resource)
				}
			}

			for path, container := range p.container {
				if len(container) == 0 {
					p.mtx.Lock()
					delete(p.container, path)
					p.mtx.Unlock()
				}
			}

			log.Debug().Msg("End cleaning")
		}
	}()
}

// Del all cache files for source file
func (p *CacheProvider) Del(resource *image.Resource) error {
	if fsutil.ContainsDotDot(resource.Path) {
		return errors.New("Invalid URL path")
	}

	path := strings.TrimPrefix(resource.Path, "/")

	container, ok := p.container[path]
	if !ok {
		return nil
	}

	size := 0

	for _, item := range container {
		size += item.Size
	}

	log.Debug().Msgf("Remove file %s", path)

	p.mtx.Lock()
	delete(p.container, path)
	p.mtx.Unlock()

	atomic.AddUint64(&p.size, ^uint64(size))

	return nil
}

// Get cached file
func (p *CacheProvider) Get(resource *image.Resource) (*image.Resource, error) {
	if fsutil.ContainsDotDot(resource.Path) {
		return nil, ErrInvalidPath
	}

	path := strings.TrimPrefix(resource.Path, "/")
	key := resource.Options.Hash()

	p.mtx.RLock()
	defer p.mtx.RUnlock()

	container, ok := p.container[path]
	if !ok {
		return nil, ErrCacheNotExist
	}

	file, ok := container[key]
	if !ok {
		return nil, ErrCacheNotExist
	}

	return &image.Resource{
		Path:       file.Path,
		Name:       file.Name,
		Options:    resource.Options,
		Body:       file.Body,
		Size:       file.Size,
		ModifiedAt: file.ModifiedAt,
	}, nil
}

// Set file to cache
func (p *CacheProvider) Set(resource *image.Resource) error {
	size := int(atomic.LoadUint64(&p.size))

	if (size + resource.Size) >= p.config.MemoryLimit {
		return fmt.Errorf("memory cache provider: allowed memory size of %d bytes exhausted", p.config.MemoryLimit)
	}

	path := strings.TrimPrefix(resource.Path, "/")
	key := resource.Options.Hash()
	_, name := filepath.Split(resource.Path)

	p.mtx.Lock()
	defer p.mtx.Unlock()

	if _, ok := p.container[path]; !ok {
		p.container[path] = make(map[string]*image.Resource)
	}

	res := &image.Resource{
		Path:       resource.Path,
		Name:       name,
		Body:       resource.Body,
		Size:       resource.Size,
		ModifiedAt: time.Now(),
	}

	p.container[path][key] = res

	atomic.AddUint64(&p.size, uint64(res.Size))

	log.Debug().Msgf("Write cache size in memory: %d", res.Size)

	return nil
}
