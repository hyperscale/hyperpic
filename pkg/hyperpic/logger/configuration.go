// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logger

import "github.com/rs/zerolog"

// Configuration struct.
type Configuration struct {
	LevelName string `mapstructure:"level"`
	Prefix    string
}

// Level type.
func (c Configuration) Level() zerolog.Level {
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
