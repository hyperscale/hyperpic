// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"github.com/rs/xlog"
	"log"
)

// ConfigLogger func
func ConfigLogger(config *LoggerConfiguration) {
	level, err := xlog.LevelFromString(config.Level)
	if err != nil {
		xlog.Fatal(err)
	}

	o := xlog.Config{
		Fields: xlog.F{
			"app": config.Prefix,
		},
		Level: level,
	}

	logger := xlog.New(o)

	log.SetOutput(logger)
	xlog.SetLogger(logger)
}
