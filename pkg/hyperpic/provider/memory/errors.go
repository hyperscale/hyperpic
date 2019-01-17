// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package memory

import "errors"

// fs errors
var (
	ErrInvalidPath    = errors.New("invalid URL path")
	ErrNotFile        = errors.New("is not a file")
	ErrCacheNotExist  = errors.New("file does not exist in cache")
	ErrSourceNotExist = errors.New("file does not exist in source")
)
