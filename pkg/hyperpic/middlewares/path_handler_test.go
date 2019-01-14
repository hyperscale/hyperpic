// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package middlewares

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/justinas/alice"
	"github.com/stretchr/testify/assert"
)

func TestPathHandlerWithBadPath(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "OK")
	}

	req := httptest.NewRequest("GET", "http://example.com/../foo.jpg?w=420&q=85&dpr=1", nil)

	w := httptest.NewRecorder()

	middleware := alice.New(
		NewPathHandler(),
	)

	middleware.ThenFunc(handler).ServeHTTP(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPathHandler(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "OK")
	}

	req := httptest.NewRequest("GET", "http://example.com/foo.jpg?w=420&q=85&dpr=1", nil)

	w := httptest.NewRecorder()

	middleware := alice.New(
		NewPathHandler(),
	)

	middleware.ThenFunc(handler).ServeHTTP(w, req)

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, []byte("OK"), body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
