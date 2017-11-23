// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestConfiguration(t *testing.T) {
	options := viper.New()

	options.SetDefault("server.port", 8080)
	options.SetDefault("logger.level", "info")

	c := NewConfiguration(options)

	assert.Equal(t, "info", c.Logger.LevelName)
	assert.Equal(t, 8080, c.Server.Port)
}
