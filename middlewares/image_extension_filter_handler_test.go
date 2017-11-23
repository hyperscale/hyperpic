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

	"github.com/hyperscale/hyperpic/config"
	"github.com/justinas/alice"
	"github.com/stretchr/testify/assert"
)

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
		NewImageExtensionFilterHandler(&config.Configuration{
			Image: &config.ImageConfiguration{
				Support: &config.ImageSupportConfiguration{
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
