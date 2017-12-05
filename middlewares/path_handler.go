// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package middlewares

import (
	"fmt"
	"net/http"

	server "github.com/euskadi31/go-server"
	"github.com/hyperscale/hyperpic/fsutil"
	"github.com/rs/zerolog/log"
)

// NewPathHandler parse query string
func NewPathHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if fsutil.ContainsDotDot(r.URL.Path) {
				// Too many programs use r.URL.Path to construct the argument to
				// serveFile. Reject the request under the assumption that happened
				// here and ".." may not be wanted.
				// Note that name might not contain "..", for example if code (still
				// incorrectly) used filepath.Join(myDir, r.URL.Path).

				log.Error().Msgf("Invalid URL path: %s", r.URL.Path)

				server.Failure(w, http.StatusBadRequest, server.ErrorMessage{
					Code:    http.StatusBadRequest,
					Message: fmt.Sprintf("Invalid URL path: %s", r.URL.Path),
				})

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
