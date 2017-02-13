package main

import (
	"math"
	"net/url"
	"strconv"
	"strings"

	bimg "gopkg.in/h2non/bimg.v1"
)

var allowedParams = map[string]string{
	"or":    "orientation",
	"w":     "int",
	"h":     "int",
	"fit":   "fit",
	"dpr":   "int",
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
		DPR:         params["dpr"].(int),
		Quality:     params["q"].(int),
		Format:      params["fm"].(bimg.ImageType),
		Orientation: params["or"].(bimg.Angle),
		Fit:         params["fit"].(FitType),
		Crop:        params["crop"].(CropType),
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
