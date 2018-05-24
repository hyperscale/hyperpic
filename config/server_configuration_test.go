// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerConfigurationAddr(t *testing.T) {
	c := &ServerConfiguration{
		Host: "127.0.0.1",
		Port: 8080,
	}

	assert.Equal(t, "127.0.0.1:8080", c.Addr())
}
