// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

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

func (s Server) imageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)

		return
	}

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

	// fetch from cache
	if resource, ok := s.cache.Get(r.URL.Path, options); ok {
		w.Header().Set("X-Image-From", "cache")

		ServeImage(w, r, resource)

		return
	}

	resource, err := s.source.Get(r.URL.Path)
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

	resource.Options = options

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
		NewParamsHandler(),
		NewClientHintsHandler(),
	)

	s.mux.HandleFunc("/health", s.healthHandler)
	s.mux.Handle("/", middleware.ThenFunc(s.imageHandler))

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      s.mux,
	}

	srv.Addr = fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)

	xlog.Infof("Server running on %s", srv.Addr)

	xlog.Fatal(srv.ListenAndServe())
}
