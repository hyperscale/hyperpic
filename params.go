// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"

	"regexp"

	"github.com/h2non/bimg"
)

var (
	colorRGBMatcher = regexp.MustCompile("^(\\d+),(\\d+),(\\d+)$")
)

var allowedParams = map[string]string{
	"or":    "orientation",
	"w":     "int",
	"h":     "int",
	"fit":   "fit",
	"dpr":   "float",
	"bri":   "int",
	"con":   "int",
	"gam":   "float",
	"sharp": "int",
	"blur":  "int",
	"pixel": "int",
	"bg":    "color",
	"q":     "int",
	"fm":    "fm",
	"crop":  "crop",
	// "filt":  "filter",
}

func ParseParams(query url.Values) *ImageOptions {
	params := make(map[string]interface{})

	for key, kind := range allowedParams {
		param := query.Get(key)
		params[key] = parseParam(param, kind)
	}

	return mapImageParams(params)
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
	case "orientation":
		return parseOrientation(param)
	case "fit":
		return parseFit(param)
	case "fm":
		return parseFormat(param)
	case "crop":
		return parseCrop(param)
	default:
		return param
	}
}

func mapImageParams(params map[string]interface{}) *ImageOptions {
	return &ImageOptions{
		Width:       params["w"].(int),
		Height:      params["h"].(int),
		DPR:         params["dpr"].(float64),
		Quality:     params["q"].(int),
		Format:      params["fm"].(bimg.ImageType),
		Orientation: params["or"].(bimg.Angle),
		Fit:         params["fit"].(FitType),
		Crop:        params["crop"].(CropType),
		Background:  params["bg"].([]uint8),
		Brightness:  params["bri"].(int),
		Contrast:    params["con"].(int),
		Gamma:       params["gam"].(float64),
		Sharpen:     params["sharp"].(int),
		Blur:        params["blur"].(int),
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

	if val == "" {
		return buf
	}

	// retrun color by name
	if color, ok := colorsToRGB[val]; ok {
		return color
	}

	if strings.Contains(val, ",") {
		for _, num := range strings.Split(val, ",") {
			n, _ := strconv.ParseUint(strings.Trim(num, " "), 10, 8)
			buf = append(buf, uint8(math.Min(float64(n), max)))
		}
	}

	if length := len(val); length == 6 || length == 3 {
		format := "%02x%02x%02x"
		factor := 1.0 / 255.0

		if length == 3 {
			format = "%1x%1x%1x"
			factor = 1.0 / 15.0
		}

		var r, g, b uint8
		n, err := fmt.Sscanf(val, format, &r, &g, &b)
		if err != nil {
			return buf
		}

		if n != 3 {
			return buf
		}

		buf = append(buf, uint8(math.Min(float64(r)*factor, max)))
		buf = append(buf, uint8(math.Min(float64(g)*factor, max)))
		buf = append(buf, uint8(math.Min(float64(b)*factor, max)))
	}

	return buf
}

func parseOrientation(val string) bimg.Angle {
	or, ok := orientationToType[val]
	if ok == false {
		return bimg.D0
	}

	return or
}

func parseFit(val string) FitType {
	fit, ok := fitToType[val]
	if ok == false {
		return FitContain
	}

	return fit
}

func parseFormat(val string) bimg.ImageType {
	return ImageType(val)
}

func parseCrop(val string) CropType {
	crop := strings.Split(val, ",")

	if len(crop) != 4 {
		return CropType{
			Width:  -1,
			Height: -1,
			X:      -1,
			Y:      -1,
		}
	}

	w, err := strconv.Atoi(crop[0])
	if err != nil {
		w = -1
	}

	h, err := strconv.Atoi(crop[1])
	if err != nil {
		h = -1
	}

	x, err := strconv.Atoi(crop[2])
	if err != nil {
		x = -1
	}

	y, err := strconv.Atoi(crop[3])
	if err != nil {
		y = -1
	}

	return CropType{
		Width:  w,
		Height: h,
		X:      x,
		Y:      y,
	}
}
