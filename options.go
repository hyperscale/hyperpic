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
	"regexp"
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
	Orientation bimg.Angle
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
			matches := cropMatcher.FindAllStringSubmatch(params["fit"].(string), 2)
			if len(matches) != 1 {
				return ImageOptions{}, fmt.Errorf("Bad fit type: %s", params["fit"].(string))
			}

			fit = FitCropFocalPoint

			x, _ := strconv.Atoi(matches[0][1])
			y, _ := strconv.Atoi(matches[0][2])

			opts.Crop = CropOptions{
				X:      x,
				Y:      y,
				Height: -1,
				Width:  -1,
			}
		}

		opts.Fit = fit
	}

	if params["w"].(int) > 0 {
		opts.Width = params["w"].(int)
	}

	if params["h"].(int) > 0 {
		opts.Height = params["h"].(int)
	}

	if params["q"].(int) > 0 {
		opts.Quality = params["q"].(int)
	}

	if params["crop"].(string) != "" {
		crop := strings.Split(params["crop"].(string), ",")

		w, err := strconv.Atoi(crop[0])
		if err != nil {
			return ImageOptions{}, fmt.Errorf("Bad crop value")
		}

		h, err := strconv.Atoi(crop[1])
		if err != nil {
			return ImageOptions{}, fmt.Errorf("Bad crop value")
		}

		x, err := strconv.Atoi(crop[2])
		if err != nil {
			return ImageOptions{}, fmt.Errorf("Bad crop value")
		}

		y, err := strconv.Atoi(crop[3])
		if err != nil {
			return ImageOptions{}, fmt.Errorf("Bad crop value")
		}

		opts.Crop = CropOptions{
			Width:  w,
			Height: h,
			X:      x,
			Y:      y,
		}
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
		Rotate: o.Orientation,
		// Interlace:    true,
		// Interpolator: bimg.Bilinear,
		NoProfile: true,
		Embed:     true,
		Enlarge:   true,
		Quality:   o.Quality,
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
		Type:           ImageType(o.Type),*/
	}

	/*if len(o.Background) != 0 {
		opts.Background = bimg.Color{o.Background[0], o.Background[1], o.Background[2]}
	}*/

	return opts
}
