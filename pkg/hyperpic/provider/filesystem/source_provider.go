// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package filesystem

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hyperscale/hyperpic/pkg/hyperpic/fsutil"
	"github.com/hyperscale/hyperpic/pkg/hyperpic/image"
	"github.com/rs/zerolog/log"
)

// SourceProvider struct.
type SourceProvider struct {
	path string
}

// NewSourceProvider func.
func NewSourceProvider(cfg *SourceConfiguration) *SourceProvider {
	return &SourceProvider{
		path: cfg.Path,
	}
}

// Set resource to file system.
func (p SourceProvider) Set(resource *image.Resource) error {
	path := p.path + "/" + strings.TrimPrefix(resource.Path, "/")

	dir, _ := filepath.Split(path)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}

	defer file.Close()

	n, err := file.Write(resource.Body)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	log.Debug().Msgf("Write source file size: %d", n)

	return nil
}

// Get resource from file system.
func (p SourceProvider) Get(resource *image.Resource) (*image.Resource, error) {
	// nolint: wsl
	if fsutil.ContainsDotDot(resource.Path) {
		// Too many programs use r.URL.Path to construct the argument to
		// serveFile. Reject the request under the assumption that happened
		// here and ".." may not be wanted.
		// Note that name might not contain "..", for example if code (still
		// incorrectly) used filepath.Join(myDir, r.URL.Path).

		return nil, ErrInvalidPath
	}

	path := p.path + "/" + strings.TrimPrefix(resource.Path, "/")

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open source file: %w", err)
	}
	defer f.Close()

	d, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat source file: %w", err)
	}

	if d.IsDir() {
		log.Debug().Msgf("%s is not a file", resource.Path)

		return nil, ErrNotFile
	}

	body, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("read source file: %w", err)
	}

	_, name := filepath.Split(resource.Path)

	return &image.Resource{
		Path:       resource.Path,
		Options:    resource.Options,
		Name:       name,
		Body:       body,
		Size:       uint64(d.Size()),
		ModifiedAt: d.ModTime(),
	}, nil
}

// Del source files.
func (p SourceProvider) Del(resource *image.Resource) error {
	if fsutil.ContainsDotDot(resource.Path) {
		return ErrInvalidPath
	}

	path := p.path + "/" + strings.TrimPrefix(resource.Path, "/")

	return os.Remove(path) // nolint: wrapcheck
}
