// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"github.com/spf13/viper"
)

// Configuration struct
type Configuration struct {
	Logger *LoggerConfiguration
	Server *ServerConfiguration
	Image  *ImageConfiguration
	Auth   *AuthConfiguration
	Doc    *DocConfiguration
}

// NewConfiguration constructor
func NewConfiguration(options *viper.Viper) *Configuration {
	return &Configuration{
		Server: &ServerConfiguration{
			Host:              options.GetString("server.host"),
			Port:              options.GetInt("server.port"),
			Debug:             options.GetBool("server.debug"),
			ShutdownTimeout:   options.GetDuration("server.shutdown_timeout"),
			WriteTimeout:      options.GetDuration("server.write_timeout"),
			ReadTimeout:       options.GetDuration("server.read_timeout"),
			ReadHeaderTimeout: options.GetDuration("server.read_header_timeout"),
		},
		Logger: &LoggerConfiguration{
			LevelName: options.GetString("logger.level"),
			Prefix:    options.GetString("logger.prefix"),
		},
		Image: &ImageConfiguration{
			Source: &ImageSourceConfiguration{
				MaxSize:  int64(options.GetInt("image.source.max_size")),
				Provider: options.GetString("image.source.provider"),
				FS: &SourceFSConfiguration{
					Path: options.GetString("image.source.fs.path"),
				},
			},
			Cache: &ImageCacheConfiguration{
				Provider: options.GetString("image.cache.provider"),
				FS: &CacheFSConfiguration{
					Path:          options.GetString("image.cache.fs.path"),
					LifeTime:      options.GetDuration("image.cache.fs.life_time"),
					CleanInterval: options.GetDuration("image.cache.fs.clean_interval"),
				},
			},
			Support: &ImageSupportConfiguration{
				Extensions: options.GetStringMap("image.support.extensions"),
			},
		},
		Auth: &AuthConfiguration{
			Secret: options.GetString("auth.secret"),
		},
		Doc: &DocConfiguration{
			Enable: options.GetBool("doc.enable"),
		},
	}
}
