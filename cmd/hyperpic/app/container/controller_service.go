// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package container

import (
	service "github.com/euskadi31/go-service"
	"github.com/hyperscale/hyperpic/cmd/hyperpic/app/config"
	"github.com/hyperscale/hyperpic/cmd/hyperpic/app/controller"
	"github.com/hyperscale/hyperpic/pkg/hyperpic/image"
	"github.com/hyperscale/hyperpic/pkg/hyperpic/provider"
)

// Services keys
const (
	DocControllerKey   = "service.controller.doc"
	ImageControllerKey = "service.controller.image"
)

func init() {
	service.Set(DocControllerKey, func(c service.Container) interface{} {
		return controller.NewDocController()
	})

	service.Set(ImageControllerKey, func(c service.Container) interface{} {
		cfg := c.Get(ConfigKey).(*config.Configuration)
		optionParser := c.Get(OptionParserKey).(*image.OptionParser)
		sourceProvider := c.Get(SourceProviderKey).(provider.SourceProvider)
		cacheProvider := c.Get(CacheProviderKey).(provider.CacheProvider)

		return controller.NewImageController(
			cfg,
			optionParser,
			sourceProvider,
			cacheProvider,
		)
	})
}
