// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/euskadi31/go-server"
	"github.com/hyperscale/hyperpic/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const applicationName = "hyperpic"

func main() {
	_ = container.Get(ServiceLoggerKey).(zerolog.Logger)
	cfg := container.Get(ServiceConfigKey).(*config.Configuration)
	router := container.Get(ServiceRouterKey).(*server.Router)

	addr := cfg.Server.Addr()

	srv := &http.Server{
		Handler:           router,
		Addr:              addr,
		WriteTimeout:      cfg.Server.WriteTimeout,
		ReadTimeout:       cfg.Server.ReadTimeout,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
	}

	go func() {
		log.Info().Msgf("Server running on %s", addr)

		if err := srv.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Msg("HTTP Server")
		}
	}()

	stop := make(chan os.Signal, 1)
	defer close(stop)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	log.Info().Msgf("Shutdown with timeout: %v", cfg.Server.ShutdownTimeout)

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("HTTP Server")
	} else {
		log.Info().Msg("Server stopped")
	}
}
