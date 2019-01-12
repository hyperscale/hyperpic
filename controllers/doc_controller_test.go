// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	server "github.com/euskadi31/go-server"
	"github.com/stretchr/testify/assert"
)

func TestDocControllerGetSwagger(t *testing.T) {
	controller, _ := NewDocController()

	router := server.NewRouter()

	router.AddController(controller)

	req := httptest.NewRequest(http.MethodGet, "/docs/swagger.yaml", nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/x-yaml", resp.Header.Get("Content-Type"))
}

func TestDocControllerGetDocs(t *testing.T) {
	controller, _ := NewDocController()

	router := server.NewRouter()

	router.AddController(controller)

	req := httptest.NewRequest(http.MethodGet, "/docs/", nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "text/html; charset=utf-8", resp.Header.Get("Content-Type"))
}
