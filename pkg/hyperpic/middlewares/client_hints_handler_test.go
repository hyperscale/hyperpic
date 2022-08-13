// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package middlewares

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hyperscale/hyperpic/pkg/hyperpic/image"
	"github.com/justinas/alice"
	"github.com/stretchr/testify/assert"
)

func TestClientHintsHandlerWithoutParamsInContext(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}

	req := httptest.NewRequest("GET", "http://example.com/foo.jpg?w=420&q=85&dpr=1", nil)

	w := httptest.NewRecorder()

	middleware := alice.New(
		NewClientHintsHandler(),
	)

	middleware.ThenFunc(handler).ServeHTTP(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestClientHintsHandler(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		options, err := OptionsFromContext(r.Context())
		assert.NoError(t, err)

		assert.Equal(t, 2.0, options.DPR)
		assert.Equal(t, 320, options.Width)
		assert.Equal(t, 65, options.Quality)

		io.WriteString(w, "OK")
	}

	req := httptest.NewRequest("GET", "http://example.com/foo.jpg?w=420&q=85&dpr=1", nil)
	req.Header.Set("DPR", "2")
	req.Header.Set("Width", "320")
	req.Header.Set("Save-Data", "on")

	w := httptest.NewRecorder()

	middleware := alice.New(
		NewOptionsHandler(image.NewOptionParser()),
		NewClientHintsHandler(),
	)

	middleware.ThenFunc(handler).ServeHTTP(w, req)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, []byte("OK"), body)
	assert.Equal(t, "2.0", resp.Header.Get("Content-DPR"))
	assert.Equal(t, []string{"DPR", "Width", "Save-Data"}, resp.Header["Vary"])
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
