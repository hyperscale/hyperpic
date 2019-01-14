// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package image

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/h2non/bimg"
	"github.com/rs/zerolog/log"
)

// FitType type
type FitType int

// Fit
const (
	FitContain FitType = iota
	FitMax
	FitFill
	FitStretch
	FitCropTopLeft
	FitCropTop
	FitCropTopRight
	FitCropLeft
	FitCropCenter
	FitCropRight
	FitCropBottomLeft
	FitCropBottom
	FitCropBottomRight
	FitCropFocalPoint
)

// FilterType type
type FilterType int

// Filter
const (
	FilterGreyscale FilterType = iota
	FilterSepia
)

var orientationToType = map[string]bimg.Angle{
	"0":   bimg.D0,
	"45":  bimg.D45,
	"90":  bimg.D90,
	"135": bimg.D135,
	"180": bimg.D180,
	"235": bimg.D235,
	"270": bimg.D270,
	"315": bimg.D315,
}

var formatToType = map[string]bimg.ImageType{
	"jpeg": bimg.JPEG,
	"jpg":  bimg.JPEG,
	"png":  bimg.PNG,
	"webp": bimg.WEBP,
	"tiff": bimg.TIFF,
}

var fitToType = map[string]FitType{
	"contain":           FitContain,
	"max":               FitMax,
	"fill":              FitFill,
	"stretch":           FitStretch,
	"crop":              FitCropCenter,
	"crop-top-left":     FitCropTopLeft,
	"crop-top":          FitCropTop,
	"crop-top-right":    FitCropTopRight,
	"crop-left":         FitCropLeft,
	"crop-center":       FitCropCenter,
	"crop-right":        FitCropRight,
	"crop-bottom-left":  FitCropBottomLeft,
	"crop-bottom":       FitCropBottom,
	"crop-bottom-right": FitCropBottomRight,
	"crop-focal-point":  FitCropFocalPoint,
}

/*
var filterToType = map[string]FilterType{
	"greyscale": FilterGreyscale,
	"sepia":     FilterSepia,
}
*/

// CropType struct
type CropType struct {
	Width  int
	Height int
	X      int
	Y      int
}

// Options represent all the supported image transformation params as first level members
type Options struct {
	Orientation bimg.Angle     `schema:"or"`
	Crop        CropType       `schema:"crop"`
	Width       int            `schema:"w"`
	Height      int            `schema:"h"`
	Fit         FitType        `schema:"fit"`
	DPR         float64        `schema:"dpr"`
	Brightness  int            `schema:"bri"`
	Contrast    int            `schema:"con"`
	Gamma       float64        `schema:"gam"`
	Sharpen     int            `schema:"sharp"`
	Blur        int            `schema:"blur"`
	Filter      FilterType     `schema:"-"`
	Background  []uint8        `schema:"bg"`
	Quality     int            `schema:"q"`
	Format      bimg.ImageType `schema:"fm"`
	Compression int            `schema:"-"`
	hash        string         `schema:"-"`
	//pixel       int            `schema:"-"`
}

// Hash return hash of options
func (o *Options) Hash() string {
	if o.hash != "" {
		return o.hash
	}

	hasher := sha256.New()
	_, _ = hasher.Write([]byte(fmt.Sprintf(
		"w=%d&h=%d&fit=%d&q=%d&fm=%d&dpr=%f&or=%d&bg=%v&bri=%d&con=%d&gam=%f&sharp=%d&blur=%d",
		o.Width,
		o.Height,
		o.Fit,
		o.Quality,
		o.Format,
		o.DPR,
		o.Orientation,
		o.Background,
		o.Brightness,
		o.Contrast,
		o.Gamma,
		o.Sharpen,
		o.Blur,
	)))
	o.hash = hex.EncodeToString(hasher.Sum(nil))

	return o.hash
}

// ToBimg creates a new bimg compatible options struct mapping the fields properly
func (o Options) ToBimg() bimg.Options {
	dpr := o.DPR

	if dpr == 0.0 {
		dpr = 1.0
	}

	width := float64(o.Width)
	height := float64(o.Height)

	if width > 0 {
		width = (width * dpr)
	}

	if height > 0 {
		height = (height * dpr)
	}

	opts := bimg.Options{
		Width:         int(width),
		Height:        int(height),
		Crop:          o.Fit == FitCropCenter,
		Rotate:        o.Orientation,
		NoProfile:     true,
		StripMetadata: true,
		Embed:         o.Fit == FitFill,
		Enlarge:       true,
		Quality:       o.Quality,
		Type:          o.Format,
		// Interlace:    true,
		// Interpolator: bimg.Bilinear,
	}

	if o.Blur > 0 {
		opts.GaussianBlur.Sigma = 0.5 * float64(o.Blur)
		opts.GaussianBlur.MinAmpl = 1.0 * float64(o.Blur)
	}

	if o.Crop.Height > 0 && o.Crop.Width > 0 {
		opts.AreaHeight = o.Crop.Height
		opts.AreaWidth = o.Crop.Width
		opts.Left = o.Crop.X
		opts.Top = o.Crop.Y
		opts.Crop = true
	}

	if len(o.Background) == 3 {
		opts.Background = bimg.Color{
			R: o.Background[0],
			G: o.Background[1],
			B: o.Background[2],
		}
	}

	if o.Fit == FitCropFocalPoint {
		opts.Crop = true
		opts.Gravity = bimg.GravitySmart
	}

	log.Debug().Msgf("options bimg: %#v", opts)

	return opts
}
