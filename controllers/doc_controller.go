// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controllers

import (
	"net/http"

	"github.com/euskadi31/go-server"
	"github.com/hyperscale/hyperpic/asset"
	"github.com/hyperscale/hyperpic/memfs"
)

// DocController struct
type DocController struct {
}

// NewDocController func
func NewDocController() (*DocController, error) {
	return &DocController{}, nil
}

// Mount endpoints
func (c DocController) Mount(r *server.Router) {
	r.AddRouteFunc("/docs/swagger.yaml", c.GetSwaggerHandler).Methods(http.MethodGet)
	r.AddRouteFunc("/docs/", c.GetDocHandler).Methods(http.MethodGet)
}

// GetSwaggerHandler endpoint
func (c DocController) GetSwaggerHandler(w http.ResponseWriter, r *http.Request) {
	name := "docs/swagger.yaml"

	w.Header().Set("Content-Type", "application/x-yaml")

	body, err := asset.Asset(name)
	if err != nil {
		server.FailureFromError(w, 0, err)

		return
	}

	info, err := asset.AssetInfo(name)
	if err != nil {
		server.FailureFromError(w, 0, err)

		return
	}

	http.ServeContent(
		w,
		r,
		info.Name(),
		info.ModTime(),
		memfs.NewBuffer(&body),
	)
}

// GetDocHandler endpoint
func (c DocController) GetDocHandler(w http.ResponseWriter, r *http.Request) {
	name := "docs/index.html"

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	body, err := asset.Asset(name)
	if err != nil {
		server.FailureFromError(w, 0, err)

		return
	}

	info, err := asset.AssetInfo(name)
	if err != nil {
		server.FailureFromError(w, 0, err)

		return
	}

	http.ServeContent(
		w,
		r,
		info.Name(),
		info.ModTime(),
		memfs.NewBuffer(&body),
	)
}
