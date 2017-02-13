package main

import (
	"context"
	"errors"
	"net/http"

	"fmt"

	"github.com/euskadi31/image-service/httputil"
	"gopkg.in/h2non/bimg.v1"
)

var (
	errContextIsNull     = errors.New("The context is null")
	errNotFountInContext = errors.New("The entry is not found in context")
)

type key int

const (
	paramsKey key = iota
)

// NewParamsContext func
func NewParamsContext(ctx context.Context, params *ImageOptions) context.Context {
	return context.WithValue(ctx, paramsKey, params)
}

// ParamsFromContext gets the params out of the context.
func ParamsFromContext(ctx context.Context) (*ImageOptions, error) {
	if ctx == nil {
		return nil, errContextIsNull
	}

	params, ok := ctx.Value(paramsKey).(*ImageOptions)
	if !ok {
		return nil, errNotFountInContext
	}

	return params, nil
}

// NewParamsHandler parse query string
func NewParamsHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			params := ParseParams(r.URL.Query())

			r = r.WithContext(NewParamsContext(r.Context(), params))

			next.ServeHTTP(w, r)
		})
	}
}

// NewClientHintsHandler parse query string
// see: http://httpwg.org/http-extensions/client-hints.html
func NewClientHintsHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			params, err := ParamsFromContext(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			if params.Format != bimg.UNKNOWN {
				next.ServeHTTP(w, r)

				return
			}

			w.Header().Set("Accept-CH", "DPR, Width")
			w.Header().Add("Vary", "Accept")

			mime := httputil.NegotiateContentType(r, []string{
				"image/webp",
				"image/jpeg",
				"image/jpg",
				"image/tiff",
				"image/png",
			}, "image/jpg")

			format := ExtractImageTypeFromMime(mime)

			if !bimg.IsTypeNameSupported(format) {
				http.Error(w, fmt.Sprintf("Format not supported"), http.StatusUnsupportedMediaType)

				return
			}

			params.Format = ImageType(format)
			w.Header().Set("Content-Type", mime)

			if dpr := r.Header.Get("DPR"); dpr != "" {
				params.DPR = parseInt(dpr)

				w.Header().Set("Content-DPR", fmt.Sprintf("%d", params.DPR))
				w.Header().Add("Vary", "DPR")
			}

			if width := r.Header.Get("Width"); width != "" {
				params.Width = parseInt(width)
				w.Header().Add("Vary", "Width")
			}

			r = r.WithContext(NewParamsContext(r.Context(), params))

			next.ServeHTTP(w, r)
		})
	}
}
