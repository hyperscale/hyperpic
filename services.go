// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/euskadi31/go-server"
	"github.com/euskadi31/go-service"
	"github.com/hyperscale/hyperpic/config"
	"github.com/hyperscale/hyperpic/controllers"
	"github.com/hyperscale/hyperpic/image"
	"github.com/hyperscale/hyperpic/provider"
	"github.com/hyperscale/hyperpic/provider/filesystem"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Service Container
var container = service.New()

var (
	basePath = "/var/lib/" + applicationName
)

// const of service name
const (
	ServiceLoggerKey          string = "service.logger"
	ServiceConfigKey                 = "service.config"
	ServiceRouterKey                 = "service.router"
	ServiceOptionParserKey           = "service.options.parser"
	ServiceSourceProviderKey         = "service.source.provider"
	ServiceCacheProviderKey          = "service.cache.provider"
	ServiceImageControllerKey        = "service.image.controller"
	ServiceDocControllerKey          = "service.doc.controller"
)

func init() {
	// Logger Service
	container.Set(ServiceLoggerKey, func(c *service.Container) interface{} {
		cfg := c.Get(ServiceConfigKey).(*config.Configuration)

		logger := zerolog.New(os.Stdout).With().
			Timestamp().
			Str("role", cfg.Logger.Prefix).
			//Str("host", host).
			Logger()

		zerolog.SetGlobalLevel(cfg.Logger.Level())

		fi, _ := os.Stdin.Stat()
		if (fi.Mode() & os.ModeCharDevice) != 0 {
			logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		}

		stdlog.SetFlags(0)
		stdlog.SetOutput(logger)

		log.Logger = logger

		return logger
	})

	// Config Service
	container.Set(ServiceConfigKey, func(c *service.Container) interface{} {
		var cfgFile string
		cmd := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		cmd.StringVar(&cfgFile, "config", "", "config file (default is $HOME/config.yaml)")

		// Ignore errors; cmd is set for ExitOnError.
		cmd.Parse(os.Args[1:])

		options := viper.New()

		if cfgFile != "" { // enable ability to specify config file via flag
			options.SetConfigFile(cfgFile)
		}

		options.SetDefault("server.host", "")
		options.SetDefault("server.port", 8080)
		options.SetDefault("server.shutdown_timeout", 10*time.Second)
		options.SetDefault("server.write_timeout", 10*time.Second)
		options.SetDefault("server.read_timeout", 10*time.Second)
		options.SetDefault("server.read_header_timeout", 10*time.Millisecond)
		options.SetDefault("logger.level", "info")
		options.SetDefault("logger.prefix", applicationName)
		options.SetDefault("auth.secret", "")
		options.SetDefault("image.source.provider", "fs")
		options.SetDefault("image.source.fs.path", basePath+"/source")
		options.SetDefault("image.cache.provider", "fs")
		options.SetDefault("image.cache.fs.path", basePath+"/cache")
		options.SetDefault("image.cache.fs.life_time", "1w")
		options.SetDefault("image.cache.fs.clean_interval", "1h")
		options.SetDefault("image.support.extensions", map[string]bool{
			"jpg":  true,
			"jpeg": true,
			"webp": true,
			"png":  true,
			"tiff": true,
		})
		options.SetDefault("doc.enable", true)

		options.SetConfigName("config") // name of config file (without extension)

		options.AddConfigPath("/etc/" + applicationName + "/")   // path to look for the config file in
		options.AddConfigPath("$HOME/." + applicationName + "/") // call multiple times to add many search paths
		options.AddConfigPath(".")

		if port := os.Getenv("PORT"); port != "" {
			os.Setenv("HYPERPIC_SERVER_PORT", port)
		}

		options.SetEnvPrefix("HYPERPIC")
		options.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		options.AutomaticEnv() // read in environment variables that match

		// If a config file is found, read it in.
		if err := options.ReadInConfig(); err == nil {
			log.Info().Msgf("Using config file: %s", options.ConfigFileUsed())
		}

		return config.NewConfiguration(options)
	})

	// Source Provider Service
	container.Set(ServiceSourceProviderKey, func(c *service.Container) interface{} {
		cfg := c.Get(ServiceConfigKey).(*config.Configuration)

		var source provider.SourceProvider

		switch cfg.Image.Source.Provider {
		case "fs":
			source = filesystem.NewSourceProvider(cfg.Image.Source.FS)
		default:
			log.Fatal().Err(fmt.Errorf("The source %s provider is not supported", cfg.Image.Source.Provider)).Msg("Source Provider")
		}

		return source
	})

	// Options Parser Service
	container.Set(ServiceOptionParserKey, func(c *service.Container) interface{} {
		return image.NewOptionParser()
	})

	// Cache Provider Service
	container.Set(ServiceCacheProviderKey, func(c *service.Container) interface{} {
		cfg := c.Get(ServiceConfigKey).(*config.Configuration)

		var cache provider.CacheProvider

		switch cfg.Image.Cache.Provider {
		case "fs":
			cache = filesystem.NewCacheProvider(cfg.Image.Cache.FS)
		default:
			log.Fatal().Err(fmt.Errorf("The cache %s provider is not supported", cfg.Image.Cache.Provider)).Msg("Cache Provider")
		}

		return cache
	})

	// Doc Controller Service
	container.Set(ServiceDocControllerKey, func(c *service.Container) interface{} {
		controller, err := controllers.NewDocController()
		if err != nil {
			log.Fatal().Err(err).Msg("Doc Controller")
		}

		return controller
	})

	// Image Controller Service
	container.Set(ServiceImageControllerKey, func(c *service.Container) interface{} {
		cfg := c.Get(ServiceConfigKey).(*config.Configuration)
		optionParser := c.Get(ServiceOptionParserKey).(*image.OptionParser)
		sourceProvider := c.Get(ServiceSourceProviderKey).(provider.SourceProvider)
		cacheProvider := c.Get(ServiceCacheProviderKey).(provider.CacheProvider)

		controller, err := controllers.NewImageController(
			cfg,
			optionParser,
			sourceProvider,
			cacheProvider,
		)
		if err != nil {
			log.Fatal().Err(err).Msg("Image Controller")
		}

		return controller
	})

	// Router Service
	container.Set(ServiceRouterKey, func(c *service.Container) interface{} {
		logger := c.Get(ServiceLoggerKey).(zerolog.Logger)
		cfg := c.Get(ServiceConfigKey).(*config.Configuration)
		imageController := c.Get(ServiceImageControllerKey).(*controllers.ImageController)
		docController := c.Get(ServiceDocControllerKey).(*controllers.DocController)

		corsHandler := cors.New(cors.Options{
			AllowCredentials: false,
			AllowedOrigins:   []string{"*"},
			AllowedMethods: []string{
				http.MethodGet,
				http.MethodOptions,
				http.MethodPost,
				http.MethodPut,
				http.MethodDelete,
			},
			AllowedHeaders: []string{
				"Authorization",
				"Content-Type",
			},
			Debug: cfg.Server.Debug,
		})

		router := server.NewRouter()

		router.Use(hlog.NewHandler(logger))
		router.Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Info().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Int("status", status).
				Int("size", size).
				Dur("duration", duration).
				Msg(fmt.Sprintf("%s %s", r.Method, r.URL.String()))
		}))
		router.Use(hlog.RemoteAddrHandler("ip"))
		router.Use(hlog.UserAgentHandler("user_agent"))
		router.Use(hlog.RefererHandler("referer"))
		router.Use(hlog.RequestIDHandler("req_id", "Request-Id"))
		router.Use(corsHandler.Handler)

		router.EnableHealthCheck()

		if cfg.Doc.Enable {
			router.AddController(docController)
		}

		router.AddController(imageController)

		return router
	})
}
