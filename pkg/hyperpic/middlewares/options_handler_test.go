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

func TestOptionsFromContext(t *testing.T) {
	opts, err := OptionsFromContext(nil)
	assert.Nil(t, opts)
	assert.EqualError(t, err, "The context is null")
}

func TestOptionsHandlerWithError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		options, err := OptionsFromContext(r.Context())
		assert.NoError(t, err)

		assert.Equal(t, 1.0, options.DPR)
		assert.Equal(t, 420, options.Width)
		assert.Equal(t, 85, options.Quality)

		io.WriteString(w, "OK")
	}

	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo.jpg?w=zeer&q=85&dpr=1", nil)

	w := httptest.NewRecorder()

	middleware := alice.New(
		NewOptionsHandler(image.NewOptionParser()),
	)

	middleware.ThenFunc(handler).ServeHTTP(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestOptionsHandler(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		options, err := OptionsFromContext(r.Context())
		assert.NoError(t, err)

		assert.Equal(t, 1.0, options.DPR)
		assert.Equal(t, 420, options.Width)
		assert.Equal(t, 85, options.Quality)

		io.WriteString(w, "OK")
	}

	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo.jpg?w=420&q=85&dpr=1", nil)

	w := httptest.NewRecorder()

	middleware := alice.New(
		NewOptionsHandler(image.NewOptionParser()),
	)

	middleware.ThenFunc(handler).ServeHTTP(w, req)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, []byte("OK"), body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
