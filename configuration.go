// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"strings"
)

type ImageConfiguration struct {
	Path string
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
			Host: "",
			Port: 8080,
		},
		Image: &ImageConfiguration{
			Path: "/var/lib/image-service",
		},
	}
}

// SourcePath return the absolute path of source folder
func (c ImageConfiguration) SourcePath() string {
	return strings.TrimSuffix(c.Path, "/") + "/source"
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
