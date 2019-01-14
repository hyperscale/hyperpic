// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package middlewares

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/hyperscale/hyperpic/cmd/hyperpic/app/config"
)

// NewAuthHandler authenticate by key
func NewAuthHandler(cfg *config.AuthConfiguration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
			if len(s) != 2 {
				http.Error(w, "Not authorized", http.StatusUnauthorized)

				return
			}

			if s[0] != "Bearer" {
				http.Error(w, "Not authorized", http.StatusUnauthorized)

				return
			}

			if subtle.ConstantTimeCompare([]byte(s[1]), []byte(cfg.Secret)) == 0 {
				http.Error(w, "Not authorized", http.StatusUnauthorized)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
