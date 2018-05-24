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

	"github.com/hyperscale/hyperpic/image"
	"github.com/justinas/alice"
	"github.com/stretchr/testify/assert"
)

func TestOptionsHandler(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		options, err := OptionsFromContext(r.Context())
		assert.NoError(t, err)

		assert.Equal(t, 1.0, options.DPR)
		assert.Equal(t, 420, options.Width)
		assert.Equal(t, 85, options.Quality)

		io.WriteString(w, "OK")
	}

	req := httptest.NewRequest("GET", "http://example.com/foo.jpg?w=420&q=85&dpr=1", nil)

	w := httptest.NewRecorder()

	middleware := alice.New(
		NewOptionsHandler(image.NewOptionParser()),
	)

	middleware.ThenFunc(handler).ServeHTTP(w, req)

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, []byte("OK"), body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
