// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"strings"
	"time"

	"github.com/rs/xlog"
	"github.com/spf13/viper"
)

var (
	basePath = "/var/lib/" + AppName
)

func init() {
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.prefix", AppName)
	viper.SetDefault("server.host", "")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("auth.secret", "")
	viper.SetDefault("image.source.provider", "fs")
	viper.SetDefault("image.source.fs.path", basePath+"/source")
	viper.SetDefault("image.cache.provider", "fs")
	viper.SetDefault("image.cache.fs.path", basePath+"/cache")
	viper.SetDefault("image.cache.fs.life_time", "1w")
	viper.SetDefault("image.cache.fs.clean_interval", "1h")
	viper.SetDefault("image.support.extensions", map[string]bool{
		"jpg":  true,
		"jpeg": true,
		"webp": true,
		"png":  true,
		"tiff": true,
	})
	viper.SetDefault("doc.enable", true)

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/" + AppName + "/")   // path to look for the config file in
	viper.AddConfigPath("$HOME/." + AppName + "/") // call multiple times to add many search paths
	viper.AddConfigPath(".")

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		xlog.Info(err)
	}

	if port := os.Getenv("PORT"); port != "" {
		os.Setenv("HYPERPIC_SERVER_PORT", port)
	}

	viper.SetEnvPrefix("HYPERPIC")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

// SourceFSConfiguration struct
type SourceFSConfiguration struct {
	Path string
}

// ImageSourceConfiguration struct
type ImageSourceConfiguration struct {
	Provider string
	FS       *SourceFSConfiguration
}

// CacheFSConfiguration struct
type CacheFSConfiguration struct {
	Path          string
	LifeTime      time.Duration
	CleanInterval time.Duration
}

// ImageCacheConfiguration struct
type ImageCacheConfiguration struct {
	Provider string
	FS       *CacheFSConfiguration
}

// ImageSupportConfiguration struct
type ImageSupportConfiguration struct {
	Extensions map[string]interface{}
}

// IsExtSupported return true if ext is supported
func (c ImageSupportConfiguration) IsExtSupported(ext string) bool {
	enable, ok := c.Extensions[ext]

	return (ok && enable.(bool))
}

// ImageConfiguration struct
type ImageConfiguration struct {
	Source  *ImageSourceConfiguration
	Cache   *ImageCacheConfiguration
	Support *ImageSupportConfiguration
}

// ServerConfiguration struct
type ServerConfiguration struct {
	Host string
	Port int
}

// LoggerConfiguration struct
type LoggerConfiguration struct {
	Level  string
	Prefix string
}

// AuthConfiguration struct
type AuthConfiguration struct {
	Secret string
}

// DocConfiguration struct
type DocConfiguration struct {
	Enable bool
}

// Configuration struct
type Configuration struct {
	Logger *LoggerConfiguration
	Server *ServerConfiguration
	Image  *ImageConfiguration
	Auth   *AuthConfiguration
	Doc    *DocConfiguration
}

// NewConfiguration constructor
func NewConfiguration() *Configuration {
	return &Configuration{
		Logger: &LoggerConfiguration{
			Level:  viper.GetString("logger.level"),
			Prefix: viper.GetString("logger.prefix"),
		},
		Server: &ServerConfiguration{
			Host: viper.GetString("server.host"),
			Port: viper.GetInt("server.port"),
		},
		Image: &ImageConfiguration{
			Source: &ImageSourceConfiguration{
				Provider: viper.GetString("image.source.provider"),
				FS: &SourceFSConfiguration{
					Path: viper.GetString("image.source.fs.path"),
				},
			},
			Cache: &ImageCacheConfiguration{
				Provider: viper.GetString("image.cache.provider"),
				FS: &CacheFSConfiguration{
					Path:          viper.GetString("image.cache.fs.path"),
					LifeTime:      viper.GetDuration("image.cache.fs.life_time"),
					CleanInterval: viper.GetDuration("image.cache.fs.clean_interval"),
				},
			},
			Support: &ImageSupportConfiguration{
				Extensions: viper.GetStringMap("image.support.extensions"),
			},
		},
		Auth: &AuthConfiguration{
			Secret: viper.GetString("auth.secret"),
		},
		Doc: &DocConfiguration{
			Enable: viper.GetBool("doc.enable"),
		},
	}
}
