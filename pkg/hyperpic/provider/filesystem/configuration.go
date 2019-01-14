// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package filesystem

import "time"

// CacheConfiguration struct
type CacheConfiguration struct {
	Path          string
	LifeTime      time.Duration `mapstructure:"life_time"`
	CleanInterval time.Duration `mapstructure:"clean_interval"`
}

// SourceConfiguration struct
type SourceConfiguration struct {
	Path string
}
