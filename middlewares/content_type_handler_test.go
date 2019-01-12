// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package middlewares

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hyperscale/hyperpic/image"
	"github.com/justinas/alice"
	"github.com/stretchr/testify/assert"
)

func TestContentTypeHandlerWithoutOptionsParserInContext(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}

	middleware := alice.New(
		NewContentTypeHandler(),
	)

	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo.jpg?w=420&q=85&dpr=1&fm=jpg", nil)

	w := httptest.NewRecorder()

	middleware.ThenFunc(handler).ServeHTTP(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestContentTypeHandler(t *testing.T) {
	tests := []struct {
		url                 string
		accept              string
		expectedBody        string
		expectedCode        int
		expectedContentType string
	}{
		{
			url:                 "http://example.com/foo.jpg?w=420&q=85&dpr=1&fm=jpg",
			accept:              "image/webp",
			expectedBody:        "Format: 1",
			expectedCode:        http.StatusOK,
			expectedContentType: "image/jpeg",
		},
		{
			url:                 "http://example.com/foo.jpg?w=420&q=85&dpr=1",
			accept:              "image/webp",
			expectedBody:        "Format: 2",
			expectedCode:        http.StatusOK,
			expectedContentType: "image/webp",
		},
		{
			url:                 "http://example.com/foo.jpg?w=420&q=85&dpr=1",
			accept:              "image/png",
			expectedBody:        "Format: 3",
			expectedCode:        http.StatusOK,
			expectedContentType: "image/png",
		},
		{
			url:                 "http://example.com/foo.jpg?w=420&q=85&dpr=1",
			accept:              "*/*",
			expectedBody:        "Format: 1",
			expectedCode:        http.StatusOK,
			expectedContentType: "image/jpeg",
		},
		{
			url:                 "http://example.com/foo.jpg?w=420&q=85&dpr=1",
			accept:              "image/flif",
			expectedBody:        "Format: 1",
			expectedCode:        http.StatusOK,
			expectedContentType: "image/jpeg",
		},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		options, err := OptionsFromContext(r.Context())
		assert.NoError(t, err)

		io.WriteString(w, fmt.Sprintf("Format: %v", options.Format))
	}

	middleware := alice.New(
		NewOptionsHandler(image.NewOptionParser()),
		NewContentTypeHandler(),
	)

	for i, tc := range tests {
		req := httptest.NewRequest("GET", tc.url, nil)
		req.Header.Set("Accept", tc.accept)

		w := httptest.NewRecorder()

		middleware.ThenFunc(handler).ServeHTTP(w, req)

		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		assert.Equal(t, tc.expectedBody, string(body))
		assert.Equal(t, tc.expectedCode, resp.StatusCode)
		assert.Equalf(t, tc.expectedContentType, resp.Header.Get("Content-Type"), "Not equal at %d", i)
	}
}
