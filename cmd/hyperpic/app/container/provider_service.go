// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package container

import (
	"fmt"

	service "github.com/euskadi31/go-service"
	"github.com/hyperscale/hyperpic/cmd/hyperpic/app/config"
	"github.com/hyperscale/hyperpic/pkg/hyperpic/provider"
	"github.com/hyperscale/hyperpic/pkg/hyperpic/provider/filesystem"
	"github.com/rs/zerolog/log"
)

// Services keys
const (
	CacheProviderKey  = "service.provider.cache"
	SourceProviderKey = "service.provider.source"
)

func init() {
	service.Set(CacheProviderKey, func(c service.Container) interface{} {
		cfg := c.Get(ConfigKey).(*config.Configuration)

		var cache provider.CacheProvider

		switch cfg.Image.Cache.Provider {
		case "fs":
			cache = filesystem.NewCacheProvider(cfg.Image.Cache.FS)
		default:
			log.Fatal().Err(fmt.Errorf("The cache %s provider is not supported", cfg.Image.Cache.Provider)).Msg("Cache Provider")
		}

		return cache
	})

	service.Set(SourceProviderKey, func(c service.Container) interface{} {
		cfg := c.Get(ConfigKey).(*config.Configuration)

		var source provider.SourceProvider

		switch cfg.Image.Source.Provider {
		case "fs":
			source = filesystem.NewSourceProvider(cfg.Image.Source.FS)
		default:
			log.Fatal().Err(fmt.Errorf("The source %s provider is not supported", cfg.Image.Source.Provider)).Msg("Source Provider")
		}

		return source
	})
}
