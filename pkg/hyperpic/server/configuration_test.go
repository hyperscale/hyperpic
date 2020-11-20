// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigurationToConfig(t *testing.T) {
	c := &Configuration{
		HTTP: &HTTPConfiguration{
			Host: "localhost",
			Port: 80,
		},
		HTTPS: &HTTPSConfiguration{
			Host: "localhost",
			Port: 443,
		},
	}

	cfg := c.ToConfig()

	assert.NotNil(t, cfg.HTTP)
	assert.NotNil(t, cfg.HTTPS)

	assert.Equal(t, "localhost", cfg.HTTP.Host)
	assert.Equal(t, 80, cfg.HTTP.Port)

	assert.Equal(t, "localhost", cfg.HTTPS.Host)
	assert.Equal(t, 443, cfg.HTTPS.Port)
}
