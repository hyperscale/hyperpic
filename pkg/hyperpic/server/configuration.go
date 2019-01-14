// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package server

import (
	"crypto/tls"
	"time"

	server "github.com/euskadi31/go-server"
)

// HTTPConfiguration struct
type HTTPConfiguration struct {
	Host string
	Port int
}

// HTTPSConfiguration struct
type HTTPSConfiguration struct {
	Host      string
	Port      int
	TLSConfig *tls.Config
	CertFile  string `mapstructure:"cert_file"`
	KeyFile   string `mapstructure:"key_file"`
}

// Configuration struct
type Configuration struct {
	HTTP              *HTTPConfiguration
	HTTPS             *HTTPSConfiguration
	ShutdownTimeout   time.Duration `mapstructure:"shutdown_timeout"`
	WriteTimeout      time.Duration `mapstructure:"write_timeout"`
	ReadTimeout       time.Duration `mapstructure:"read_timeout"`
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
	IdleTimeout       time.Duration `mapstructure:"idle_timeout"`
	Profiling         bool
	Metrics           bool
	HealthCheck       bool
}

// ToConfig returns server.Configuration
func (c Configuration) ToConfig() *server.Configuration {
	cfg := &server.Configuration{
		ShutdownTimeout:   c.ShutdownTimeout,
		WriteTimeout:      c.WriteTimeout,
		ReadTimeout:       c.ReadTimeout,
		ReadHeaderTimeout: c.ReadHeaderTimeout,
		IdleTimeout:       c.IdleTimeout,
		Profiling:         c.Profiling,
		Metrics:           c.Metrics,
		HealthCheck:       c.HealthCheck,
	}

	if c.HTTP != nil {
		cfg.HTTP = &server.HTTPConfiguration{
			Host: c.HTTP.Host,
			Port: c.HTTP.Port,
		}
	}

	if c.HTTPS != nil {
		cfg.HTTPS = &server.HTTPSConfiguration{
			Host:      c.HTTPS.Host,
			Port:      c.HTTPS.Port,
			TLSConfig: c.HTTPS.TLSConfig,
			CertFile:  c.HTTPS.CertFile,
			KeyFile:   c.HTTPS.KeyFile,
		}
	}

	return cfg
}
