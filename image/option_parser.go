// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package image

import (
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
}

func (p OptionParser) angleConverter(s string) reflect.Value {
	value, ok := orientationToType[s]
	if ok == false {
		return reflect.ValueOf(bimg.D0)
	}

	return reflect.ValueOf(value)
}

func (p OptionParser) formatConverter(s string) reflect.Value {
	value, ok := formatToType[s]
	if ok == false {
		return reflect.ValueOf(bimg.UNKNOWN)
	}

	return reflect.ValueOf(value)
}

func (p OptionParser) fitConverter(s string) reflect.Value {
	value, ok := fitToType[s]
	if ok == false {
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
