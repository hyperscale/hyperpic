// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"github.com/hyperscale/hyperpic/pkg/hyperpic/provider/filesystem"
	"github.com/hyperscale/hyperpic/pkg/hyperpic/provider/memory"
)

// ImageCacheConfiguration struct
type ImageCacheConfiguration struct {
	Provider string
	FS       *filesystem.CacheConfiguration
	Memory   *memory.CacheConfiguration
}
