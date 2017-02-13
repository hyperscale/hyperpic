package main

import (
	"crypto/md5"
	"encoding/hex"

	"regexp"

	"fmt"

	bimg "gopkg.in/h2non/bimg.v1"
)

var (
	cropMatcher = regexp.MustCompile("^crop-(\\d+)-(\\d+)")
)

type FitType int
type FilterType int

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
}

var filterToType = map[string]FilterType{
	"greyscale": FilterGreyscale,
	"sepia":     FilterSepia,
}

type CropType struct {
	Width  int
	Height int
	X      int
	Y      int
}

// ImageOptions represent all the supported image transformation params as first level members
type ImageOptions struct {
	Orientation bimg.Angle
	Crop        CropType
	Width       int
	Height      int
	Fit         FitType
	DPR         int
	Brightness  int
	Contrast    int
	Gamma       float64
	Sharpen     int
	Blur        int
	pixel       int
	Filter      FilterType
	Background  bimg.Color
	Quality     int
	Format      bimg.ImageType
	Compression int
	hash        string
}

// Hash return hash of options
func (o *ImageOptions) Hash() string {
	if o.hash != "" {
		return o.hash
	}

	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf(
		"w=%d&h=%d&fit=%d&q=%d&fm=%d&dpr=%d&or=%d",
		o.Width,
		o.Height,
		o.Fit,
		o.Quality,
		o.Format,
		o.DPR,
		o.Orientation,
	)))
	o.hash = hex.EncodeToString(hasher.Sum(nil))

	return o.hash
}

// BimgOptions creates a new bimg compatible options struct mapping the fields properly
func BimgOptions(o *ImageOptions) bimg.Options {
	opts := bimg.Options{
		Width:     o.Width,
		Height:    o.Height,
		Crop:      o.Fit == FitCropCenter,
		Rotate:    o.Orientation,
		NoProfile: true,
		Embed:     true,
		Enlarge:   true,
		Quality:   o.Quality,
		Type:      o.Format,
		// Interlace:    true,
		// Interpolator: bimg.Bilinear,
	}

	/*if len(o.Background) != 0 {
		opts.Background = bimg.Color{o.Background[0], o.Background[1], o.Background[2]}
	}*/

	return opts
}
