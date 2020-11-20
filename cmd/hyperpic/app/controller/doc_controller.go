// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controller

import (
	"net/http"

	server "github.com/euskadi31/go-server"
	"github.com/hyperscale/hyperpic/cmd/hyperpic/app/asset"
	"github.com/hyperscale/hyperpic/pkg/hyperpic/memfs"
)

type docController struct {
}

// NewDocController func
func NewDocController() server.Controller {
	return &docController{}
}

// Mount endpoints
func (c docController) Mount(r *server.Router) {
	r.AddRouteFunc("/docs/swagger.yaml", c.getSwaggerHandler).Methods(http.MethodGet)
	r.AddRouteFunc("/docs/", c.getDocHandler).Methods(http.MethodGet)
}

// GET /docs/swagger.yaml
func (c docController) getSwaggerHandler(w http.ResponseWriter, r *http.Request) {
	name := "docs/swagger.yaml"

	w.Header().Set("Content-Type", "application/x-yaml")

	body, _ := asset.Asset(name)

	info, _ := asset.AssetInfo(name)

	http.ServeContent(
		w,
		r,
		info.Name(),
		info.ModTime(),
		memfs.NewBuffer(&body),
	)
}

// GET /docs/
func (c docController) getDocHandler(w http.ResponseWriter, r *http.Request) {
	name := "docs/index.html"

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	body, _ := asset.Asset(name)

	info, _ := asset.AssetInfo(name)

	http.ServeContent(
		w,
		r,
		info.Name(),
		info.ModTime(),
		memfs.NewBuffer(&body),
	)
}
