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

	"github.com/hyperscale/hyperpic/cmd/hyperpic/app/config"
	"github.com/justinas/alice"
	"github.com/stretchr/testify/assert"
)

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
		NewAuthHandler(&config.AuthConfiguration{
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
