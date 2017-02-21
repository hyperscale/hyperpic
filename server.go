// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"os"
	"time"

	filetype "gopkg.in/h2non/filetype.v0"

	"io/ioutil"

	"github.com/euskadi31/image-service/httputil"
	"github.com/justinas/alice"
	"github.com/rs/xlog"
)

// Server struct
type Server struct {
	mux    *http.ServeMux
	config *Configuration
	source SourceProvider
	cache  CacheProvider
}

// NewServer constructor
func NewServer(config *Configuration) *Server {
	var source SourceProvider
	var cache CacheProvider

	switch config.Image.Source.Provider {
	case "fs":
		source = NewFileSystemSourceProvider(config.Image.Source.FS)
	}

	switch config.Image.Cache.Provider {
	case "fs":
		cache = NewFileSystemCacheProvider(config.Image.Cache.FS)
	}

	return &Server{
		mux:    http.NewServeMux(),
		config: config,
		source: source,
		cache:  cache,
	}
}

// SetSourceProvider to server
func (s *Server) SetSourceProvider(source SourceProvider) {
	s.source = source
}

// SetCacheProvider to server
func (s *Server) SetCacheProvider(cache CacheProvider) {
	s.cache = cache
}

func (s Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)

		return
	}

	w.Write([]byte("OK"))
}

func (s Server) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found", http.StatusNotFound)
}

func (s Server) deleteHandler(w http.ResponseWriter, r *http.Request) {
	if containsDotDot(r.URL.Path) {
		// Too many programs use r.URL.Path to construct the argument to
		// serveFile. Reject the request under the assumption that happened
		// here and ".." may not be wanted.
		// Note that name might not contain "..", for example if code (still
		// incorrectly) used filepath.Join(myDir, r.URL.Path).

		httputil.Failure(w, http.StatusBadRequest, httputil.ErrorMessage{
			Code:    0,
			Message: "Invalid URL path",
		})

		return
	}

	resource := &Resource{
		Path: r.URL.Path,
	}

	response := map[string]bool{
		"cache":  false,
		"source": false,
	}

	if from := r.URL.Query().Get("from"); from != "" {
		switch from {
		case "all":
			response["cache"] = s.cache.Del(resource)
			response["source"] = s.source.Del(resource)
		case "source":
			response["source"] = s.source.Del(resource)
		case "cache":
			response["cache"] = s.cache.Del(resource)
		}
	}

	httputil.JSON(w, http.StatusOK, resource)
}

func (s Server) uploadHandler(w http.ResponseWriter, r *http.Request) {
	if containsDotDot(r.URL.Path) {
		// Too many programs use r.URL.Path to construct the argument to
		// serveFile. Reject the request under the assumption that happened
		// here and ".." may not be wanted.
		// Note that name might not contain "..", for example if code (still
		// incorrectly) used filepath.Join(myDir, r.URL.Path).

		httputil.Failure(w, http.StatusBadRequest, httputil.ErrorMessage{
			Code:    0,
			Message: "Invalid URL path",
		})

		return
	}

	resource := &Resource{
		Path: r.URL.Path,
	}

	r.ParseMultipartForm(32 << 20)

	file, _, err := r.FormFile("image")
	if err != nil {
		httputil.FailureFromError(w, http.StatusBadRequest, err)

		return
	}

	defer file.Close()

	body, err := ioutil.ReadAll(file)
	if err != nil {
		httputil.FailureFromError(w, http.StatusBadRequest, err)

		return
	}

	resource.Body = body

	if err := s.source.Set(resource); err != nil {
		httputil.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	// delete cache from source file
	go s.cache.Del(resource)

	mimeType := http.DetectContentType(body)

	// If cannot infer the type, infer it via magic numbers
	if mimeType == "application/octet-stream" {
		kind, err := filetype.Get(body)
		if err == nil && kind.MIME.Value != "" {
			mimeType = kind.MIME.Value
		}
	}

	h := md5.New()
	h.Write(body)

	httputil.JSON(w, http.StatusCreated, map[string]interface{}{
		"file": r.URL.Path,
		"size": len(body),
		"type": mimeType,
		"hash": fmt.Sprintf("%x", h.Sum(nil)),
	})
}

func (s Server) imageHandler(w http.ResponseWriter, r *http.Request) {
	if containsDotDot(r.URL.Path) {
		// Too many programs use r.URL.Path to construct the argument to
		// serveFile. Reject the request under the assumption that happened
		// here and ".." may not be wanted.
		// Note that name might not contain "..", for example if code (still
		// incorrectly) used filepath.Join(myDir, r.URL.Path).

		http.Error(w, "Invalid URL path", http.StatusBadRequest)

		return
	}

	options, err := ParamsFromContext(r.Context())
	if err != nil {
		xlog.Errorf("Error while parsing options: %v", err.Error())
		http.Error(w, "Error while parsing options", http.StatusBadRequest)

		return
	}
	// xlog.Infof("options: %#v", options)

	resource := &Resource{
		Path:    r.URL.Path,
		Options: options,
	}

	// fetch from cache
	if resource, ok := s.cache.Get(resource); ok {
		w.Header().Set("X-Image-From", "cache")

		ServeImage(w, r, resource)

		return
	}

	resource, err = s.source.Get(resource)
	if err != nil {
		if os.IsNotExist(err) {
			msg := fmt.Sprintf("File %s not found", r.URL.Path)

			xlog.Info(msg)
			http.Error(w, msg, http.StatusNotFound)

			return
		}

		xlog.Error(err)
		http.Error(w, err.Error(), http.StatusNotFound)

		return
	}

	if err := ProcessImage(resource); err != nil {
		xlog.Errorf("Error while processing the image: %v", err.Error())
		http.Error(w, "Error while processing the image", http.StatusInternalServerError)

		return
	}

	w.Header().Set("X-Image-From", "source")

	ServeImage(w, r, resource)

	// save resource in cache
	go func(r *Resource) {
		if err := s.cache.Set(r); err != nil {
			xlog.Error(err)
		}
	}(resource)
}

// ListenAndServe service
func (s *Server) ListenAndServe() {
	middleware := alice.New(
		NewLoggerHandler(),
		NewImageExtensionFilterHandler(s.config),
	)

	readMiddleware := middleware.Append(
		NewParamsHandler(),
		NewClientHintsHandler(),
	)

	authMiddleware := middleware.Append(
		NewAuthHandler(s.config.Auth),
	)

	s.mux.HandleFunc("/favicon.ico", s.notFoundHandler)
	s.mux.HandleFunc("/health", s.healthHandler)
	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			readMiddleware.ThenFunc(s.imageHandler).ServeHTTP(w, r)
		case http.MethodPost:
			authMiddleware.ThenFunc(s.uploadHandler).ServeHTTP(w, r)
		case http.MethodDelete:
			authMiddleware.ThenFunc(s.deleteHandler).ServeHTTP(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      s.mux,
	}

	srv.Addr = fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)

	xlog.Infof("Server running on %s", srv.Addr)

	xlog.Fatal(srv.ListenAndServe())
}
