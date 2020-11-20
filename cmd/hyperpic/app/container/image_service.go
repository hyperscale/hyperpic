// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package container

import (
	service "github.com/euskadi31/go-service"
	"github.com/hyperscale/hyperpic/pkg/hyperpic/image"
)

// Services keys
const (
	ImageOptionParserKey = "service.image.options.parser"
	ImageProcessorKey    = "service.image.processor"
)

func init() {
	service.Set(ImageOptionParserKey, func(c service.Container) interface{} {
		return image.NewOptionParser()
	})

	service.Set(ImageProcessorKey, func(c service.Container) interface{} {
		return image.NewProcessor()
	})
}
