// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package image

import (
	"net/http/httptest"
	"testing"

	"github.com/h2non/bimg"
	"github.com/stretchr/testify/assert"
)

func TestOptionParserParse(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8574/stock-photo-103005233.jpg?w=400&h=400&fit=crop&dpr=2&or=45&fm=webp", nil)

	parser := NewOptionParser()

	options, err := parser.Parse(req)
	assert.NoError(t, err)

	assert.Equal(t, 400, options.Width)
	assert.Equal(t, 400, options.Height)
	assert.Equal(t, FitCropCenter, options.Fit)
	assert.Equal(t, float64(2), options.DPR)
	assert.Equal(t, bimg.D45, options.Orientation)
	assert.Equal(t, bimg.WEBP, options.Format)
}
