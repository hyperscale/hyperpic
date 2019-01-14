// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

// ImageConfiguration struct
type ImageConfiguration struct {
	Source  *ImageSourceConfiguration
	Cache   *ImageCacheConfiguration
	Support *ImageSupportConfiguration
}
