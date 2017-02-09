package main

import (
	"crypto/md5"
	"encoding/hex"
	"math"
	"net/url"
	"strconv"
	"strings"

	"fmt"

	bimg "gopkg.in/h2non/bimg.v1"
)

type OrientationType int
type FitType int
type FilterType int

const (
	OrientationAuto OrientationType = 999
	Orientation0    OrientationType = 0
	Orientation90   OrientationType = 90
	Orientation180  OrientationType = 180
	Orientation270  OrientationType = 270
)

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

var orientationToType = map[string]OrientationType{
	"auto": OrientationAuto,
	"0":    Orientation0,
	"90":   Orientation90,
	"180":  Orientation180,
	"270":  Orientation270,
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

var allowedOptions = map[string]string{
	"or":     "orientation",
	"crop":   "string",
	"w":      "int",
	"h":      "int",
	"fit":    "fit",
	"dpr":    "int",
	"bri":    "int",
	"con":    "int",
	"gam":    "float",
	"sharp":  "int",
	"blur":   "int",
	"pixel":  "int",
	"filt":   "filter",
	"bg":     "color",
	"border": "",
	"q":      "int",
	"fm":     "fm",
}

type CropOptions struct {
	Width  int
	Height int
	X      int
	Y      int
}

// ImageOptions represent all the supported image transformation params as first level members
type ImageOptions struct {
	Orientation OrientationType
	Crop        CropOptions
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
	query       url.Values
}

func ParseImageOptions(query url.Values) (ImageOptions, error) {
	params := make(map[string]interface{})

	for key, kind := range allowedOptions {
		param := query.Get(key)
		params[key] = parseParam(param, kind)
	}

	opts := ImageOptions{}

	if params["or"].(string) != "" {
		or, ok := orientationToType[params["or"].(string)]
		if ok == false {
			return ImageOptions{}, fmt.Errorf("Bad orientation type: %s", params["or"].(string))
		}

		opts.Orientation = or
	}

	if params["fit"].(string) != "" {
		fit, ok := fitToType[params["fit"].(string)]
		if ok == false {
			return ImageOptions{}, fmt.Errorf("Bad fit type: %s", params["fit"].(string))
		}

		opts.Fit = fit
	}

	if params["w"].(int) > 0 {
		opts.Width = params["w"].(int)
	}

	if params["h"].(int) > 0 {
		opts.Height = params["h"].(int)
	}

	return opts, nil
	/*
		return ImageOptions{
			Orientation: or,
			// Crop:        params["crop"].(CropOptions),
			Width:      params["w"].(int),
			Height:     params["h"].(int),
			Fit:        fit,
			DPR:        params["dpr"].(int),
			Brightness: params["bri"].(int),
			Contrast:   params["con"].(int),
			Gamma:      params["gam"].(float64),
			Sharpen:    params["sharp"].(int),
			Blur:       params["blur"].(int),
			pixel:      params["pixel"].(int),
			//Filter:      params["filt"].(FilterType),
			//Background:  params["filt"].(bimg.Color),
			Quality: params["q"].(int),
			//Format:      params["fm"].(bimg.ImageType),
			Compression: params["c"].(int),
			query:       query,
		}, nil*/
}

// Hash return hash of options
func (o *ImageOptions) Hash() string {
	if o.hash != "" {
		return o.hash
	}

	hasher := md5.New()
	hasher.Write([]byte(o.query.Encode()))

	o.hash = hex.EncodeToString(hasher.Sum(nil))

	return o.hash
}

func parseParam(param, kind string) interface{} {
	switch kind {
	case "int":
		return parseInt(param)
	case "float":
		return parseFloat(param)
	case "color":
		return parseColor(param)
	case "bool":
		return parseBool(param)
	default:
		return param
	}
}

func parseBool(val string) bool {
	value, _ := strconv.ParseBool(val)
	return value
}

func parseInt(param string) int {
	return int(math.Floor(parseFloat(param) + 0.5))
}

func parseFloat(param string) float64 {
	val, _ := strconv.ParseFloat(param, 64)
	return math.Abs(val)
}

func parseColor(val string) []uint8 {
	const max float64 = 255
	buf := []uint8{}
	if val != "" {
		for _, num := range strings.Split(val, ",") {
			n, _ := strconv.ParseUint(strings.Trim(num, " "), 10, 8)
			buf = append(buf, uint8(math.Min(float64(n), max)))
		}
	}
	return buf
}

func parseGravity(val string) bimg.Gravity {
	val = strings.TrimSpace(strings.ToLower(val))

	switch val {
	case "south":
		return bimg.GravitySouth
	case "north":
		return bimg.GravityNorth
	case "east":
		return bimg.GravityEast
	case "west":
		return bimg.GravityWest
	default:
		return bimg.GravityCentre
	}
}

// BimgOptions creates a new bimg compatible options struct mapping the fields properly
func BimgOptions(o ImageOptions) bimg.Options {
	opts := bimg.Options{
		Width:  o.Width,
		Height: o.Height,
		Crop:   o.Fit == FitCropCenter,
		/*Flip:           o.Flip,
		Flop:           o.Flop,
		Quality:        o.Quality,
		Compression:    o.Compression,
		NoAutoRotate:   o.NoRotation,
		NoProfile:      o.NoProfile,
		Force:          o.Force,
		Gravity:        o.Gravity,
		Embed:          o.Embed,
		Extend:         o.Extend,
		Interpretation: o.Colorspace,
		Type:           ImageType(o.Type),
		Rotate:         bimg.Angle(o.Rotate),*/
	}

	/*if len(o.Background) != 0 {
		opts.Background = bimg.Color{o.Background[0], o.Background[1], o.Background[2]}
	}*/

	return opts
}
