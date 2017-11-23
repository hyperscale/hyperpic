// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package middlewares

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/hyperscale/hyperpic/config"
	"github.com/rs/zerolog/log"
)

// NewImageExtensionFilterHandler filtered url for accept image file only
func NewImageExtensionFilterHandler(cfg *config.Configuration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ext := strings.ToLower(filepath.Ext(r.URL.Path))
			ext = ext[1:]

			if cfg.Image.Support.IsExtSupported(ext) == false {
				msg := fmt.Sprintf("File %s is not supported", r.URL.Path)

				log.Debug().Msg(msg)

				http.Error(w, msg, http.StatusNotFound)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
