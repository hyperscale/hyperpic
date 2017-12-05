// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package image

import (
	"net/url"
	"testing"

	"github.com/h2non/bimg"
	"github.com/stretchr/testify/assert"
)

func TestParseParams(t *testing.T) {
	u, err := url.Parse("http://localhost:8574/stock-photo-103005233.jpg?w=400&h=400&fit=crop&dpr=2&or=45&fm=jpg")
	assert.NoError(t, err)

	params := ParseParams(u.Query())

	assert.Equal(t, 400, params.Width)
	assert.Equal(t, 400, params.Height)
	assert.Equal(t, FitCropCenter, params.Fit)
	assert.Equal(t, float64(2), params.DPR)
	assert.Equal(t, bimg.D45, params.Orientation)
	assert.Equal(t, bimg.JPEG, params.Format)
}
