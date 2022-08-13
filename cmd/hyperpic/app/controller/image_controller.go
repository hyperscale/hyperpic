// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controller

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"

	server "github.com/euskadi31/go-server"
	"github.com/euskadi31/go-server/response"
	"github.com/h2non/filetype"
	"github.com/hyperscale/hyperpic/cmd/hyperpic/app/config"
	"github.com/hyperscale/hyperpic/cmd/hyperpic/app/metrics"
	"github.com/hyperscale/hyperpic/pkg/hyperpic/httputil"
	"github.com/hyperscale/hyperpic/pkg/hyperpic/image"
	"github.com/hyperscale/hyperpic/pkg/hyperpic/middlewares"
	"github.com/hyperscale/hyperpic/pkg/hyperpic/provider"
	"github.com/justinas/alice"
	"github.com/rs/zerolog/hlog"
)

type imageController struct {
	cfg            *config.Configuration
	optionParser   *image.OptionParser
	imageProcessor image.Processor
	sourceProvider provider.SourceProvider
	cacheProvider  provider.CacheProvider
}

// NewImageController func
func NewImageController(
	cfg *config.Configuration,
	optionParser *image.OptionParser,
	imageProcessor image.Processor,
	sourceProvider provider.SourceProvider,
	cacheProvider provider.CacheProvider,
) server.Controller {
	return &imageController{
		cfg:            cfg,
		optionParser:   optionParser,
		imageProcessor: imageProcessor,
		sourceProvider: sourceProvider,
		cacheProvider:  cacheProvider,
	}
}

// Mount endpoints
func (c imageController) Mount(r *server.Router) {
	chain := alice.New(
		middlewares.NewPathHandler(),
		middlewares.NewImageExtensionFilterHandler(c.cfg),
	)

	public := chain.Append(
		middlewares.NewOptionsHandler(c.optionParser),
		middlewares.NewContentTypeHandler(),
		middlewares.NewClientHintsHandler(),
	)

	private := chain.Append(
		middlewares.NewAuthHandler(c.cfg.Auth),
	)

	r.AddPrefixRoute("/", public.ThenFunc(c.getHandler)).Methods(http.MethodGet)
	r.AddPrefixRoute("/", private.ThenFunc(c.postHandler)).Methods(http.MethodPost)
	r.AddPrefixRoute("/", private.ThenFunc(c.deleteHandler)).Methods(http.MethodDelete)
}

// GET /:file
func (c imageController) getHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := hlog.FromRequest(r)

	options, err := middlewares.OptionsFromContext(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error while parsing options")

		http.Error(w, "Error while parsing options", http.StatusBadRequest)

		return
	}
	// xlog.Infof("options: %#v", options)

	resource := &image.Resource{
		Path:    r.URL.Path,
		Options: options,
	}

	// w.Header().Set("Link", `</worker/client-hints.js>; rel="serviceworker"`)

	// fetch from cache
	if resource, err := c.cacheProvider.Get(resource); err == nil {
		w.Header().Set("X-Image-From", "cache")

		httputil.ServeImage(w, r, resource)

		metrics.CacheHit.With(map[string]string{}).Add(1)
		metrics.ImageDeliveredBytes.With(map[string]string{}).Add(float64(resource.Size))

		return
	}

	resource, err = c.sourceProvider.Get(resource)
	if err != nil {
		if os.IsNotExist(err) {
			msg := fmt.Sprintf("File %s not found", r.URL.Path)

			log.Info().Msg(msg)

			http.Error(w, msg, http.StatusNotFound)

			return
		}

		log.Error().Err(err).Msg("Source Provider")

		http.Error(w, err.Error(), http.StatusNotFound)

		return
	}

	if err := c.imageProcessor.ProcessImage(resource); err != nil {
		log.Error().Err(err).Msg("Error while processing the image")

		http.Error(w, "Error while processing the image", http.StatusInternalServerError)

		return
	}

	w.Header().Set("X-Image-From", "source")

	httputil.ServeImage(w, r, resource)

	// save resource in cache
	go func(r *image.Resource) {
		if err := c.cacheProvider.Set(r); err != nil {
			log.Error().Err(err).Msg("Cache Provider")
		}
	}(resource)

	metrics.CacheMiss.With(map[string]string{}).Add(1)
	metrics.ImageDeliveredBytes.With(map[string]string{}).Add(float64(resource.Size))
}

func (c imageController) parseImageFileFromRequest(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, errors.New("missing form body")
	}

	r.Body = http.MaxBytesReader(w, r.Body, c.cfg.Image.Source.MaxSize)

	ct := r.Header.Get("Content-Type")
	// RFC 7231, section 3.1.1.5 - empty type
	//   MAY be treated as application/octet-stream
	if ct == "" {
		ct = "application/octet-stream"
	}

	ct, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return nil, err
	}

	switch ct {
	case "multipart/form-data":
		if err := r.ParseMultipartForm(c.cfg.Image.Source.MaxSize); err != nil {
			return nil, err
		}

		file, _, err := r.FormFile("image")
		if err != nil {
			return nil, err
		}

		defer file.Close()

		return io.ReadAll(file)
	default:
		return io.ReadAll(r.Body)
	}
}

// POST /:file
func (c imageController) postHandler(w http.ResponseWriter, r *http.Request) {
	log := hlog.FromRequest(r)

	resource := &image.Resource{
		Path: r.URL.Path,
	}

	body, err := c.parseImageFileFromRequest(w, r)
	if err != nil && err.Error() == "http: request body too large" {
		log.Error().Err(err).Msg("parseImageFileFromRequest failed")

		response.FailureFromError(w, http.StatusRequestEntityTooLarge, err)

		return
	} else if err != nil {
		log.Error().Err(err).Msg("parseImageFileFromRequest failed")

		response.FailureFromError(w, http.StatusBadRequest, err)

		return
	}

	resource.Body = body

	if err := c.sourceProvider.Set(resource); err != nil {
		response.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	// delete cache from source file
	go func() {
		if err := c.cacheProvider.Del(resource); err != nil {
			log.Error().Err(err).Msg("CacheProvider.Del failed")
		}
	}()

	mimeType := http.DetectContentType(body)

	// If cannot infer the type, infer it via magic numbers
	if mimeType == "application/octet-stream" {
		kind, err := filetype.Get(body)
		if err == nil && kind.MIME.Value != "" {
			mimeType = kind.MIME.Value
		}
	}

	h := sha256.New()
	length, _ := h.Write(body)

	response.Encode(w, r, http.StatusCreated, map[string]interface{}{
		"file": r.URL.Path,
		"size": length,
		"type": mimeType,
		"hash": fmt.Sprintf("%x", h.Sum(nil)),
	})

	metrics.ImageReceivedBytes.With(map[string]string{}).Add(float64(length))
}

// DELETE /:file
func (c imageController) deleteHandler(w http.ResponseWriter, r *http.Request) {
	resource := &image.Resource{
		Path: r.URL.Path,
	}

	resp := map[string]bool{
		"cache":  false,
		"source": false,
	}

	if from := r.URL.Query().Get("from"); from != "" {
		switch from {
		case "source":
			resp["cache"] = (c.cacheProvider.Del(resource) == nil)
			resp["source"] = (c.sourceProvider.Del(resource) == nil)
		default:
			resp["cache"] = (c.cacheProvider.Del(resource) == nil)
		}
	}

	response.Encode(w, r, http.StatusOK, resp)
}
