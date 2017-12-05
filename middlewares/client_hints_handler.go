// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package middlewares

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
)

func parseInt(value string) int {
	return int(math.Floor(parseFloat(value) + 0.5))
}

func parseFloat(value string) float64 {
	val, _ := strconv.ParseFloat(value, 64)

	return math.Abs(val)
}

// NewClientHintsHandler parse query string
// see: http://httpwg.org/http-extensions/client-hints.html
func NewClientHintsHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			options, err := OptionsFromContext(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			if dpr := r.Header.Get("DPR"); dpr != "" {
				options.DPR = parseFloat(dpr)

				w.Header().Set("Content-DPR", fmt.Sprintf("%.1f", options.DPR))
				w.Header().Add("Vary", "DPR")
			}

			if width := r.Header.Get("Width"); width != "" {
				options.Width = parseInt(width)

				w.Header().Add("Vary", "Width")
			}

			if saveData := r.Header.Get("Save-Data"); saveData == "on" {
				options.Quality = 65

				w.Header().Add("Vary", "Save-Data")
			}

			r = r.WithContext(NewOptionsContext(r.Context(), options))

			next.ServeHTTP(w, r)
		})
	}
}
