// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package image

import (
	"bytes"
	"image"
	_ "image/png"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {
	o := &Options{
		Width:  20,
		Height: 20,
		Fit:    FitCropCenter,
	}

	file, err := os.Open("../../../_resources/hyperpic.png")
	assert.NoError(t, err)

	in, err := io.ReadAll(file)
	assert.NoError(t, err)

	p := &processor{}

	out, err := p.process(in, o.ToBimg())
	assert.NoError(t, err)

	assert.Equal(t, "image/png", out.Mime)

	image, _, err := image.DecodeConfig(bytes.NewReader(out.Body))
	assert.NoError(t, err)

	assert.Equal(t, 20, image.Height)
	assert.Equal(t, 20, image.Width)
}

func TestProcessImage(t *testing.T) {
	o := &Options{
		Width:  20,
		Height: 20,
		Fit:    FitCropCenter,
	}

	file, err := os.Open("../../../_resources/hyperpic.png")
	assert.NoError(t, err)

	in, err := io.ReadAll(file)
	assert.NoError(t, err)

	res := &Resource{
		Body:    in,
		Options: o,
	}

	p := NewProcessor()

	assert.NoError(t, p.ProcessImage(res))

	assert.Equal(t, "image/png", res.MimeType)

	res = &Resource{
		Body:    []byte{0x01},
		Options: o,
	}

	assert.EqualError(t, p.ProcessImage(res), "MimeType application/octet-stream is not supported")
}
