// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"os"
	"os/signal"
	"syscall"

	server "github.com/euskadi31/go-server"
	service "github.com/euskadi31/go-service"
	"github.com/hyperscale/hyperpic/cmd/hyperpic/app/container"
	"github.com/rs/zerolog/log"
)

// Run Hyperpic api server
func Run() error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	_ = service.Get(container.LoggerKey)

	router := service.Get(container.RouterKey).(*server.Server)

	log.Info().Msg("Rinning")

	go func() {
		log.Info().Msg("Rinning HTTP Router")
		if err := router.Run(); err != nil {
			log.Error().Err(err).Msg("server.Run() failed")
		}
	}()

	<-sig

	log.Info().Msg("Shutdown")

	return router.Shutdown()
}
