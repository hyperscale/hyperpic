// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"github.com/hyperscale/hyperpic/pkg/hyperpic/logger"
	"github.com/hyperscale/hyperpic/pkg/hyperpic/provider/filesystem"
	"github.com/hyperscale/hyperpic/pkg/hyperpic/server"
)

// Configuration struct
type Configuration struct {
	Logger *logger.Configuration
	Server *server.Configuration
	Image  *ImageConfiguration
	Auth   *AuthConfiguration
	Doc    *DocConfiguration
}

// NewConfiguration constructor
func NewConfiguration() *Configuration {
	return &Configuration{
		Server: &server.Configuration{},
		Logger: &logger.Configuration{},
		Image: &ImageConfiguration{
			Source: &ImageSourceConfiguration{
				FS: &filesystem.SourceConfiguration{},
			},
			Cache: &ImageCacheConfiguration{
				FS: &filesystem.CacheConfiguration{},
			},
			Support: &ImageSupportConfiguration{},
		},
		Auth: &AuthConfiguration{},
		Doc:  &DocConfiguration{},
	}
}
