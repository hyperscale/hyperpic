// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package container

import (
	stdlog "log"
	"os"

	service "github.com/euskadi31/go-service"
	"github.com/hyperscale/hyperpic/cmd/hyperpic/app/config"
	"github.com/hyperscale/hyperpic/pkg/hyperpic/version"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Services keys
const (
	LoggerKey = "service.logger"
)

func init() {
	service.Set(LoggerKey, func(c service.Container) interface{} {
		cfg := c.Get(ConfigKey).(*config.Configuration)

		zerolog.SetGlobalLevel(cfg.Logger.Level())

		logger := zerolog.New(os.Stdout).With().
			Timestamp().
			Str("role", cfg.Logger.Prefix).
			Str("version", version.Version).
			Logger()

		fi, err := os.Stdin.Stat()
		if err != nil {
			log.Fatal().Err(err).Msg("Stdin.Stat failed")
		}

		if (fi.Mode() & os.ModeCharDevice) != 0 {
			logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		}

		stdlog.SetFlags(0)
		stdlog.SetOutput(logger)

		log.Logger = logger

		return logger // zerolog.Logger
	})
}
