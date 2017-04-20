// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"errors"
	"net/http"

	"fmt"

	"crypto/subtle"
	"path/filepath"
	"strings"
	"time"

	"github.com/h2non/bimg"
	"github.com/hyperscale/hyperpic/httputil"
	"github.com/rs/xlog"
	"github.com/whitedevops/colors"
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

// NewLoggerHandler log request
func NewLoggerHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			xlog.Infof("  %s   %s   %s  %s", fmtDuration(start), fmtStatus(w), fmtMethod(r), fmtPath(r.URL.Path))
		})
	}
}

// NewImageExtensionFilterHandler filtered url for accept image file only
func NewImageExtensionFilterHandler(config *Configuration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ext := strings.ToLower(filepath.Ext(r.URL.Path))
			ext = ext[1:]

			if config.Image.Support.IsExtSupported(ext) == false {
				msg := fmt.Sprintf("File %s is not supported", r.URL.Path)

				xlog.Debug(msg)
				http.Error(w, msg, http.StatusNotFound)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
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

// NewAuthHandler authenticate by key
func NewAuthHandler(config *AuthConfiguration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
			if len(s) != 2 {
				http.Error(w, "Not authorized", 401)

				return
			}

			if s[0] != "Bearer" {
				http.Error(w, "Not authorized", 401)

				return
			}

			if subtle.ConstantTimeCompare([]byte(s[1]), []byte(config.Secret)) == 0 {
				http.Error(w, "Not authorized", 401)

				return
			}

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

			if params.Format == bimg.UNKNOWN {
				mime := httputil.NegotiateContentType(r, []string{
					"image/jpg",
					"image/webp",
					"image/jpeg",
					"image/tiff",
					"image/png",
				}, "image/jpg")

				format := ExtractImageTypeFromMime(mime)

				xlog.Debugf("Format extracted form mime: %s => %s", mime, format)

				if !IsFormatSupported(format) {
					http.Error(w, fmt.Sprintf("Format not supported"), http.StatusUnsupportedMediaType)

					return
				}

				params.Format = ImageType(format)
				w.Header().Set("Content-Type", mime)
				w.Header().Add("Vary", "Accept")
			}

			if dpr := r.Header.Get("DPR"); dpr != "" {
				params.DPR = parseFloat(dpr)

				w.Header().Set("Content-DPR", fmt.Sprintf("%f", params.DPR))
				w.Header().Add("Vary", "DPR")
			}

			if width := r.Header.Get("Width"); width != "" {
				params.Width = parseInt(width)
				w.Header().Add("Vary", "Width")
			}

			if saveData := r.Header.Get("Save-Data"); saveData == "on" {
				params.Quality = 65
				w.Header().Add("Vary", "Save-Data")
			}

			r = r.WithContext(NewParamsContext(r.Context(), params))

			next.ServeHTTP(w, r)
		})
	}
}

func fmtDuration(start time.Time) string {
	return fmt.Sprintf("%s%s%13s%s", colors.ResetAll, colors.Dim, time.Since(start), colors.ResetAll)
}

func fmtStatus(w http.ResponseWriter) string {
	code := httputil.ResponseStatus(w)

	color := colors.White

	switch {
	case code >= 200 && code <= 299:
		color += colors.BackgroundGreen
	case code >= 300 && code <= 399:
		color += colors.BackgroundCyan
	case code >= 400 && code <= 499:
		color += colors.BackgroundYellow
	default:
		color += colors.BackgroundRed
	}

	return fmt.Sprintf("%s%s %3d %s", colors.ResetAll, color, code, colors.ResetAll)
}

func fmtMethod(r *http.Request) string {
	var color string

	switch r.Method {
	case "GET":
		color += colors.Green
	case "POST":
		color += colors.Cyan
	case "PUT", "PATCH":
		color += colors.Blue
	case "DELETE":
		color += colors.Red
	}

	return fmt.Sprintf("%s%s%s%s", colors.ResetAll, color, r.Method, colors.ResetAll)
}

func fmtPath(path string) string {
	return fmt.Sprintf("%s%s%s%s", colors.ResetAll, colors.Dim, path, colors.ResetAll)
}
