// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package middlewares

import (
	"net/http"

	"github.com/h2non/bimg"
	"github.com/hyperscale/hyperpic/httputil"
	"github.com/hyperscale/hyperpic/image"
	"github.com/rs/zerolog/log"
)

// NewContentTypeHandler negotiate content type
func NewContentTypeHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			options, err := OptionsFromContext(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			if options.Format != bimg.UNKNOWN {
				next.ServeHTTP(w, r)

				return
			}

			mime := httputil.NegotiateContentType(r, []string{
				"image/jpg",
				"image/webp",
				"image/jpeg",
				"image/tiff",
				"image/png",
			}, "image/jpg")

			format := image.ExtractImageTypeFromMime(mime)

			log.Debug().Msgf("Format extracted form mime: %s => %s", mime, format)

			/*if !IsFormatSupported(format) {
				http.Error(w, fmt.Sprintf("Format not supported"), http.StatusUnsupportedMediaType)

				return
			}*/

			options.Format = image.ExtensionToType(format)

			w.Header().Set("Content-Type", mime)
			w.Header().Add("Vary", "Accept")

			r = r.WithContext(NewOptionsContext(r.Context(), options))

			next.ServeHTTP(w, r)
		})
	}
}
