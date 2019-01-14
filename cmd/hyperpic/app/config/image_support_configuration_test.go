// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageSupportConfigurationIsExtSupported(t *testing.T) {
	c := &ImageSupportConfiguration{
		Extensions: map[string]interface{}{
			"jpg":  true,
			"jpeg": true,
			"webp": true,
			"png":  true,
			"tiff": true,
		},
	}

	assert.True(t, c.IsExtSupported("jpg"))
	assert.True(t, c.IsExtSupported("jpeg"))
	assert.True(t, c.IsExtSupported("webp"))
	assert.True(t, c.IsExtSupported("png"))
	assert.True(t, c.IsExtSupported("tiff"))
	assert.False(t, c.IsExtSupported("svg"))
}
