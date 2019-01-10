// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

// ImageSourceConfiguration struct
type ImageSourceConfiguration struct {
	MaxSize  int64
	Provider string
	FS       *SourceFSConfiguration
}
