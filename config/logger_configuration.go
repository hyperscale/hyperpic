// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import "github.com/rs/zerolog"

// LoggerConfiguration struct
type LoggerConfiguration struct {
	LevelName string
	Prefix    string
}

// Level type
func (c LoggerConfiguration) Level() zerolog.Level {
	switch c.LevelName {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.Disabled
	}
}
