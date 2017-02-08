// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

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
	return c.Path + "/source"
}

// CachePath return the absolute path of cache folder
func (c ImageConfiguration) CachePath() string {
	return c.Path + "/cache"
}
