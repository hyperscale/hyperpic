// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfiguration(t *testing.T) {
	c := NewConfiguration()

	assert.NotNil(t, c.Logger)
	assert.NotNil(t, c.Server)
}
