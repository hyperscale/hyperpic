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

func TestBadFormatOptionParser(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8574/stock-photo-103005233.jpg?w=400&h=400&fit=bad&dpr=2&or=45&fm=bar", nil)

	parser := NewOptionParser()

	options, err := parser.Parse(req)
	assert.NoError(t, err)

	assert.Equal(t, bimg.UNKNOWN, options.Format)
}

func TestBadFitOptionParser(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8574/stock-photo-103005233.jpg?w=400&h=400&fit=bad&dpr=2&or=45&fm=jpg", nil)

	parser := NewOptionParser()

	options, err := parser.Parse(req)
	assert.NoError(t, err)

	assert.Equal(t, FitContain, options.Fit)
}

func TestBadOrientationOptionParser(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8574/stock-photo-103005233.jpg?w=400&h=400&fit=bad&dpr=2&or=42&fm=jpg", nil)

	parser := NewOptionParser()

	options, err := parser.Parse(req)
	assert.NoError(t, err)

	assert.Equal(t, bimg.D0, options.Orientation)
}

func TestBackgroundColorOptionParser(t *testing.T) {
	assertions := []struct {
		value    string
		expected []uint8
	}{
		{
			value:    "http://localhost:8574/stock-photo-103005233.jpg?bg=fff",
			expected: []uint8{255, 255, 255},
		},
		{
			value:    "http://localhost:8574/stock-photo-103005233.jpg?bg=azf",
			expected: []uint8{},
		},
		{
			value:    "http://localhost:8574/stock-photo-103005233.jpg?bg=ffffff",
			expected: []uint8{255, 255, 255},
		},
		{
			value:    "http://localhost:8574/stock-photo-103005233.jpg?bg=%23ffffff",
			expected: []uint8{255, 255, 255},
		},
		{
			value:    "http://localhost:8574/stock-photo-103005233.jpg?bg=",
			expected: []uint8{},
		},
		{
			value:    "http://localhost:8574/stock-photo-103005233.jpg?bg=255,255,255",
			expected: []uint8{255, 255, 255},
		},
		{
			value:    "http://localhost:8574/stock-photo-103005233.jpg?bg=255,255",
			expected: []uint8{},
		},
		{
			value:    "http://localhost:8574/stock-photo-103005233.jpg?bg=teal",
			expected: []uint8{0, 128, 128},
		},
	}

	for _, assertion := range assertions {
		req := httptest.NewRequest("GET", assertion.value, nil)

		parser := NewOptionParser()

		options, err := parser.Parse(req)
		assert.NoError(t, err)

		assert.Equal(t, assertion.expected, options.Background)
	}
}

func TestBadCropOptionParser(t *testing.T) {
	assertions := []struct {
		value    string
		expected CropType
	}{
		{
			value: "http://localhost:8574/stock-photo-103005233.jpg?crop=10,10,10",
			expected: CropType{
				Width:  -1,
				Height: -1,
				X:      -1,
				Y:      -1,
			},
		},
		{
			value: "http://localhost:8574/stock-photo-103005233.jpg?crop=A,11,12,13",
			expected: CropType{
				Width:  -1,
				Height: 11,
				X:      12,
				Y:      13,
			},
		},
		{
			value: "http://localhost:8574/stock-photo-103005233.jpg?crop=10,A,12,13",
			expected: CropType{
				Width:  10,
				Height: -1,
				X:      12,
				Y:      13,
			},
		},
		{
			value: "http://localhost:8574/stock-photo-103005233.jpg?crop=10,11,A,13",
			expected: CropType{
				Width:  10,
				Height: 11,
				X:      -1,
				Y:      13,
			},
		},
		{
			value: "http://localhost:8574/stock-photo-103005233.jpg?crop=10,11,12,A",
			expected: CropType{
				Width:  10,
				Height: 11,
				X:      12,
				Y:      -1,
			},
		},
	}

	for _, assertion := range assertions {
		req := httptest.NewRequest("GET", assertion.value, nil)

		parser := NewOptionParser()

		options, err := parser.Parse(req)
		assert.NoError(t, err)

		assert.Equal(t, assertion.expected.Height, options.Crop.Height)
		assert.Equal(t, assertion.expected.Width, options.Crop.Width)
		assert.Equal(t, assertion.expected.X, options.Crop.X)
		assert.Equal(t, assertion.expected.Y, options.Crop.Y)
	}
}
