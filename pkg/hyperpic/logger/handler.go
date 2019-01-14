// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logger

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

// Handler for zerolog
func Handler(r *http.Request, status, size int, duration time.Duration) {
	rlog := hlog.FromRequest(r)

	var evt *zerolog.Event

	switch {
	case status >= 200 && status <= 299:
		evt = rlog.Info()
	case status >= 300 && status <= 399:
		evt = rlog.Info()
	case status >= 400 && status <= 499:
		evt = rlog.Warn()
	default:
		evt = rlog.Error()
	}

	evt.
		Str("method", r.Method).
		Str("url", r.URL.String()).
		Int("status", status).
		Int("size", size).
		Dur("duration", duration).
		Msgf("%s %s", r.Method, r.URL.Path)
}
