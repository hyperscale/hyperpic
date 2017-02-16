// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"time"

	"strconv"

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
	return &Server{
		mux:    http.NewServeMux(),
		config: config,
		source: NewFileSystemSourceProvider(config.Image.SourcePath()),
		cache:  NewFileSystemCacheProvider(config.Image.CachePath()),
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
		HTTPFailure(w, http.StatusMethodNotAllowed)

		return
	}

	w.Write([]byte("OK"))
}

func (s Server) imageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		HTTPFailure(w, http.StatusMethodNotAllowed)

		return
	}

	if containsDotDot(r.URL.Path) {
		// Too many programs use r.URL.Path to construct the argument to
		// serveFile. Reject the request under the assumption that happened
		// here and ".." may not be wanted.
		// Note that name might not contain "..", for example if code (still
		// incorrectly) used filepath.Join(myDir, r.URL.Path).
		HTTPFailure(w, http.StatusBadRequest)

		return
	}

	options, err := ParamsFromContext(r.Context())
	if err != nil {
		xlog.Errorf("Error while parsing options: %v", err.Error())

		HTTPFailure(w, http.StatusBadRequest)

		return
	}
	xlog.Infof("options: %#v", options)
	/*if options.IsValid() == false {
		HTTPFailure(w, http.StatusBadRequest)

		return
	}*/

	// fetch from cache
	if resource, ok := s.cache.Get(r.URL.Path, options); ok {
		w.Header().Set("X-Image-From", "cache")

		http.ServeFile(w, r, resource.CachePath)

		return
	}

	resource, err := s.source.Get(r.URL.Path)
	if err != nil {
		HTTPFailure(w, http.StatusNotFound)

		return
	}

	if err := ProcessImage(resource, options); err != nil {
		xlog.Errorf("Error while processing the image: %v", err.Error())

		HTTPFailure(w, http.StatusInternalServerError)

		return
	}

	// w.Header().Set("Content-Type", resource.MimeType)
	w.Header().Set("Content-Length", strconv.Itoa(len(resource.Body)))
	w.Header().Set("X-Image-From", "source")

	w.Write(resource.Body)

	//http.ServeFile(w, r, resource.SourcePath)

	// save resource in cache
	go func(r *Resource) {
		if err := s.cache.Set(r, options); err != nil {
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
