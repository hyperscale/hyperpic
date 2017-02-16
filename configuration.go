// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"github.com/rs/xlog"
	"github.com/spf13/viper"
	"strings"
)

var (
	basePath = "/var/lib/" + AppName
)

func init() {
	viper.SetDefault("server.host", "")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("image.source.provider", "fs")
	viper.SetDefault("image.source.fs.path", basePath+"/source")
	viper.SetDefault("image.cache.provider", "fs")
	viper.SetDefault("image.cache.fs.path", basePath+"/cache")

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/" + AppName + "/")   // path to look for the config file in
	viper.AddConfigPath("$HOME/." + AppName + "/") // call multiple times to add many search paths
	viper.AddConfigPath(".")

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		xlog.Fatalf("Fatal error config file: %s", err)
	}
}

type SourceFSConfiguration struct {
	Path string
}

type ImageSourceConfiguration struct {
	Provider string
	FS       *SourceFSConfiguration
}

type CacheFSConfiguration struct {
	Path string
}

type ImageCacheConfiguration struct {
	Provider string
	FS       *CacheFSConfiguration
}

type ImageConfiguration struct {
	Source *ImageSourceConfiguration
	Cache  *ImageCacheConfiguration
}

type ServerConfiguration struct {
	Host string
	Port int
}

type Configuration struct {
	Server *ServerConfiguration
	Image  *ImageConfiguration
}

// NewConfiguration constructor
func NewConfiguration() *Configuration {
	return &Configuration{
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
					Path: viper.GetString("image.cache.fs.path"),
				},
			},
		},
	}
}

// SourcePath return the absolute path of source folder
func (c ImageConfiguration) SourcePath() string {
	return strings.TrimSuffix(c.Source.Path, "/") + "/source"
}

func (c ImageConfiguration) SourcePathWithFile(file string) string {
	return c.SourcePath() + "/" + strings.TrimPrefix(file, "/")
}

// CachePath return the absolute path of cache folder
func (c ImageConfiguration) CachePath() string {
	return strings.TrimSuffix(c.Path, "/") + "/cache"
}

func (c ImageConfiguration) CachePathWithFile(file string) string {
	return c.CachePath() + "/" + strings.TrimPrefix(file, "/")
}
