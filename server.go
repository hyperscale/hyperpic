// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/rs/xlog"
	"net/http"
	"os"
	"strings"
	"time"
)

type Server struct {
	mux    *http.ServeMux
	config *Configuration
}

// NewServer constructor
func NewServer(config *Configuration) *Server {
	return &Server{
		mux:    http.NewServeMux(),
		config: config,
	}
}

func (s Server) failure(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
}

func (s Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		s.failure(w, http.StatusMethodNotAllowed)

		return
	}

	w.Write([]byte("OK"))
}

func (s Server) imageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		s.failure(w, http.StatusMethodNotAllowed)

		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/")

	fullPath := s.config.Image.SourcePath() + "/" + path

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		s.failure(w, http.StatusNotFound)

		return
	}

	http.ServeFile(w, r, fullPath)

	//w.Write([]byte(fullPath))
}

// ListenAndServe service
func (s *Server) ListenAndServe() {
	s.mux.HandleFunc("/health", s.healthHandler)
	s.mux.HandleFunc("/", s.imageHandler)

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      s.mux,
	}

	srv.Addr = fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)

	xlog.Infof("Server running on %s", srv.Addr)

	xlog.Fatal(srv.ListenAndServe())
}
