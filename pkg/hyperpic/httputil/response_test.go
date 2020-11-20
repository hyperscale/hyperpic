// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package httputil

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hyperscale/hyperpic/pkg/hyperpic/image"
	"github.com/stretchr/testify/assert"
)

func TestServeImage(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo.jpg", nil)
	w := httptest.NewRecorder()

	ServeImage(w, req, &image.Resource{
		Name:       "foo.jpg",
		ModifiedAt: time.Now().AddDate(1955, 02, 24),
		Body:       []byte("bar"),
		Size:       3,
	})

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "image/jpeg", resp.Header.Get("Content-Type"))
	assert.Equal(t, "bar", string(body))
}
