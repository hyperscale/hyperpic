// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"fmt"

	"github.com/justinas/alice"
	"github.com/stretchr/testify/assert"
)

func TestParamsHandler(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		params, err := ParamsFromContext(r.Context())
		assert.NoError(t, err)

		assert.Equal(t, 1.0, params.DPR)
		assert.Equal(t, 420, params.Width)
		assert.Equal(t, 85, params.Quality)

		io.WriteString(w, "OK")
	}

	req := httptest.NewRequest("GET", "http://example.com/foo.jpg?w=420&q=85&dpr=1", nil)

	w := httptest.NewRecorder()

	middleware := alice.New(
		NewParamsHandler(),
	)

	middleware.ThenFunc(handler).ServeHTTP(w, req)

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, []byte("OK"), body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestImageExtensionFilterHandler(t *testing.T) {
	tests := []struct {
		url          string
		expectedBody string
		expectedCode int
	}{
		{
			url:          "http://example.com/foo.jpg?w=420&q=85&dpr=1",
			expectedBody: "File /foo.jpg is not supported\n",
			expectedCode: http.StatusNotFound,
		},
		{
			url:          "http://example.com/foo.jpeg?w=420&q=85&dpr=1",
			expectedBody: "OK",
			expectedCode: http.StatusOK,
		},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "OK")
	}

	middleware := alice.New(
		NewImageExtensionFilterHandler(&Configuration{
			Image: &ImageConfiguration{
				Support: &ImageSupportConfiguration{
					Extensions: map[string]interface{}{
						"jpg":  false,
						"jpeg": true,
						"png":  true,
						"webp": true,
					},
				},
			},
		}),
	)

	for _, tc := range tests {
		req := httptest.NewRequest("GET", tc.url, nil)

		w := httptest.NewRecorder()

		middleware.ThenFunc(handler).ServeHTTP(w, req)

		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		assert.Equal(t, tc.expectedBody, string(body))
		assert.Equal(t, tc.expectedCode, resp.StatusCode)
	}

}

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
		params, err := ParamsFromContext(r.Context())
		assert.NoError(t, err)

		assert.Equal(t, 2.0, params.DPR)
		assert.Equal(t, 320, params.Width)
		assert.Equal(t, 65, params.Quality)

		io.WriteString(w, "OK")
	}

	req := httptest.NewRequest("GET", "http://example.com/foo.jpg?w=420&q=85&dpr=1", nil)
	req.Header.Set("DPR", "2")
	req.Header.Set("Width", "320")
	req.Header.Set("Save-Data", "on")

	w := httptest.NewRecorder()

	middleware := alice.New(
		NewParamsHandler(),
		NewClientHintsHandler(),
	)

	middleware.ThenFunc(handler).ServeHTTP(w, req)

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, []byte("OK"), body)
	assert.Equal(t, "2.0", resp.Header.Get("Content-DPR"))
	assert.Equal(t, "DPR, Width, Save-Data", resp.Header.Get("Vary"))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAuthHandler(t *testing.T) {
	tests := []struct {
		authorization string
		expectedBody  string
		expectedCode  int
	}{
		{
			authorization: "bad",
			expectedBody:  "Not authorized\n",
			expectedCode:  http.StatusUnauthorized,
		},
		{
			authorization: "bad password",
			expectedBody:  "Not authorized\n",
			expectedCode:  http.StatusUnauthorized,
		},
		{
			authorization: "Bearer bad password",
			expectedBody:  "Not authorized\n",
			expectedCode:  http.StatusUnauthorized,
		},
		{
			authorization: "Bearer foo",
			expectedBody:  "OK",
			expectedCode:  http.StatusOK,
		},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "OK")
	}

	middleware := alice.New(
		NewAuthHandler(&AuthConfiguration{
			Secret: "foo",
		}),
	)

	for _, tc := range tests {
		req := httptest.NewRequest("GET", "http://example.com/foo.jpg?w=420&q=85&dpr=1", nil)
		req.Header.Set("Authorization", tc.authorization)

		w := httptest.NewRecorder()

		middleware.ThenFunc(handler).ServeHTTP(w, req)

		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		assert.Equal(t, tc.expectedBody, string(body))
		assert.Equal(t, tc.expectedCode, resp.StatusCode)
	}
}

func TestContentTypeHandler(t *testing.T) {
	tests := []struct {
		url          string
		accept       string
		expectedBody string
		expectedCode int
	}{
		{
			url:          "http://example.com/foo.jpg?w=420&q=85&dpr=1&fm=jpg",
			accept:       "image/webp",
			expectedBody: "Format: 1",
			expectedCode: http.StatusOK,
		},
		{
			url:          "http://example.com/foo.jpg?w=420&q=85&dpr=1",
			accept:       "image/webp",
			expectedBody: "Format: 2",
			expectedCode: http.StatusOK,
		},
		{
			url:          "http://example.com/foo.jpg?w=420&q=85&dpr=1",
			accept:       "image/png",
			expectedBody: "Format: 3",
			expectedCode: http.StatusOK,
		},
		{
			url:          "http://example.com/foo.jpg?w=420&q=85&dpr=1",
			accept:       "*/*",
			expectedBody: "Format: 1",
			expectedCode: http.StatusOK,
		},
		{
			url:          "http://example.com/foo.jpg?w=420&q=85&dpr=1",
			accept:       "image/flif",
			expectedBody: "Format: 1",
			expectedCode: http.StatusOK,
		},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		params, err := ParamsFromContext(r.Context())
		assert.NoError(t, err)

		io.WriteString(w, fmt.Sprintf("Format: %v", params.Format))
	}

	middleware := alice.New(
		NewParamsHandler(),
		NewContentTypeHandler(),
	)

	for _, tc := range tests {
		req := httptest.NewRequest("GET", tc.url, nil)
		req.Header.Set("Accept", tc.accept)

		w := httptest.NewRecorder()

		middleware.ThenFunc(handler).ServeHTTP(w, req)

		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		assert.Equal(t, tc.expectedBody, string(body))
		assert.Equal(t, tc.expectedCode, resp.StatusCode)
	}
}

func TestLoggerHandler(t *testing.T) {
	middleware := alice.New(
		NewLoggerHandler(),
	)

	ts := httptest.NewServer(middleware.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "OK")
	}))
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	assert.NoError(t, err)

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, []byte("OK"), body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
