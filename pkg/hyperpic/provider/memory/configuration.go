// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package memory

import "time"

// CacheConfiguration struct
type CacheConfiguration struct {
	LifeTime      time.Duration `mapstructure:"life_time"`
	CleanInterval time.Duration `mapstructure:"clean_interval"`
	MemoryLimit   uint64        `mapstructure:"memory_limit"`
}

// SourceConfiguration struct
type SourceConfiguration struct {
	MemoryLimit uint64 `mapstructure:"memory_limit"`
}
