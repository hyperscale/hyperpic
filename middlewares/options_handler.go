// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package middlewares

import (
	"context"
	"net/http"

	"github.com/euskadi31/go-server"
	"github.com/hyperscale/hyperpic/image"
	"github.com/rs/zerolog/log"
)

// NewOptionsContext func
func NewOptionsContext(ctx context.Context, options *image.Options) context.Context {
	return context.WithValue(ctx, optionsKey, options)
}

// OptionsFromContext gets the options out of the context.
func OptionsFromContext(ctx context.Context) (*image.Options, error) {
	if ctx == nil {
		return nil, errContextIsNull
	}

	options, ok := ctx.Value(optionsKey).(*image.Options)
	if !ok {
		return nil, errNotFountInContext
	}

	return options, nil
}

// NewOptionsHandler parse query string
func NewOptionsHandler(optionParser *image.OptionParser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			options, err := optionParser.Parse(r)
			if err != nil {
				log.Error().Err(err).Msg("Options Parser")

				server.FailureFromError(w, http.StatusBadRequest, err)

				return
			}

			r = r.WithContext(NewOptionsContext(r.Context(), options))

			next.ServeHTTP(w, r)
		})
	}
}
