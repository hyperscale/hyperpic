// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package image

import (
	"encoding/hex"
	"fmt"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gorilla/schema"
	"github.com/h2non/bimg"
)

// OptionParser struct
type OptionParser struct {
	decoder *schema.Decoder
}

// NewOptionParser func
func NewOptionParser() *OptionParser {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	parser := &OptionParser{
		decoder: decoder,
	}

	parser.register()

	return parser
}

func (p *OptionParser) register() {
	p.decoder.RegisterConverter(bimg.Angle(0), p.angleConverter)
	p.decoder.RegisterConverter(bimg.ImageType(0), p.formatConverter)
	p.decoder.RegisterConverter(FitType(0), p.fitConverter)
	p.decoder.RegisterConverter(CropType{}, p.cropConverter)
	p.decoder.RegisterConverter([]uint8{}, p.colorConverter)
}

func (p OptionParser) colorConverter(s string) reflect.Value {
	const max float64 = 255

	buf := []uint8{}

	if s == "" {
		return reflect.ValueOf(buf)
	}

	// retrun color by name
	if color, ok := colorsToRGB[s]; ok {
		return reflect.ValueOf(color)
	}

	if strings.Contains(s, ",") {
		parts := strings.Split(s, ",")
		if len(parts) != 3 {
			return reflect.ValueOf(buf)
		}

		for _, num := range parts {
			n, _ := strconv.ParseUint(strings.Trim(num, " "), 10, 8)
			buf = append(buf, uint8(math.Min(float64(n), max)))
		}

		return reflect.ValueOf(buf)
	}

	if strings.HasPrefix(s, "#") {
		s = strings.Replace(s, "#", "", 1)
	}

	if len(s) == 3 {
		s = fmt.Sprintf("%c%c%c%c%c%c", s[0], s[0], s[1], s[1], s[2], s[2])
	}

	d, err := hex.DecodeString(s)
	if err != nil {
		return reflect.ValueOf(buf)
	}

	buf = append(buf, uint8(d[0]))
	buf = append(buf, uint8(d[1]))
	buf = append(buf, uint8(d[2]))

	return reflect.ValueOf(buf)
}

func (p OptionParser) angleConverter(s string) reflect.Value {
	value, ok := orientationToType[s]
	if !ok {
		return reflect.ValueOf(bimg.D0)
	}

	return reflect.ValueOf(value)
}

func (p OptionParser) formatConverter(s string) reflect.Value {
	value, ok := formatToType[s]
	if !ok {
		return reflect.ValueOf(bimg.UNKNOWN)
	}

	return reflect.ValueOf(value)
}

func (p OptionParser) fitConverter(s string) reflect.Value {
	value, ok := fitToType[s]
	if !ok {
		return reflect.ValueOf(FitContain)
	}

	return reflect.ValueOf(value)
}

func (p OptionParser) cropConverter(s string) reflect.Value {
	crop := strings.Split(s, ",")

	if len(crop) != 4 {
		return reflect.ValueOf(CropType{
			Width:  -1,
			Height: -1,
			X:      -1,
			Y:      -1,
		})
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

	return reflect.ValueOf(CropType{
		Width:  w,
		Height: h,
		X:      x,
		Y:      y,
	})
}

// Parse Option from url
func (p OptionParser) Parse(r *http.Request) (*Options, error) {
	option := &Options{}

	if err := p.decoder.Decode(option, r.URL.Query()); err != nil {
		return nil, err
	}

	return option, nil
}
