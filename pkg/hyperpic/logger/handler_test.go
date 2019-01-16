// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logger

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestHandler2xx(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "einfo")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	logger := zerolog.New(tmpfile).With().Logger()

	req := httptest.NewRequest(http.MethodGet, "http://foo.com/health", nil)

	req = req.WithContext(logger.WithContext(req.Context()))

	Handler(req, http.StatusOK, 45, 2*time.Millisecond)

	actual, err := ioutil.ReadFile(tmpfile.Name())
	assert.NoError(t, err)

	assert.Equal(t, "{\"level\":\"info\",\"method\":\"GET\",\"url\":\"http://foo.com/health\",\"status\":200,\"size\":45,\"duration\":2,\"message\":\"GET /health\"}\n", string(actual))
}

func TestHandler3xx(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "einfo")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	logger := zerolog.New(tmpfile).With().Logger()

	req := httptest.NewRequest(http.MethodGet, "http://foo.com/health", nil)

	req = req.WithContext(logger.WithContext(req.Context()))

	Handler(req, http.StatusPermanentRedirect, 45, 2*time.Millisecond)

	actual, err := ioutil.ReadFile(tmpfile.Name())
	assert.NoError(t, err)

	assert.Equal(t, "{\"level\":\"info\",\"method\":\"GET\",\"url\":\"http://foo.com/health\",\"status\":308,\"size\":45,\"duration\":2,\"message\":\"GET /health\"}\n", string(actual))
}

func TestHandler4xx(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "einfo")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	logger := zerolog.New(tmpfile).With().Logger()

	req := httptest.NewRequest(http.MethodGet, "http://foo.com/health", nil)

	req = req.WithContext(logger.WithContext(req.Context()))

	Handler(req, http.StatusNotFound, 45, 2*time.Millisecond)

	actual, err := ioutil.ReadFile(tmpfile.Name())
	assert.NoError(t, err)

	assert.Equal(t, "{\"level\":\"warn\",\"method\":\"GET\",\"url\":\"http://foo.com/health\",\"status\":404,\"size\":45,\"duration\":2,\"message\":\"GET /health\"}\n", string(actual))
}

func TestHandler5xx(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "einfo")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	logger := zerolog.New(tmpfile).With().Logger()

	req := httptest.NewRequest(http.MethodGet, "http://foo.com/health", nil)

	req = req.WithContext(logger.WithContext(req.Context()))

	Handler(req, http.StatusInternalServerError, 45, 2*time.Millisecond)

	actual, err := ioutil.ReadFile(tmpfile.Name())
	assert.NoError(t, err)

	assert.Equal(t, "{\"level\":\"error\",\"method\":\"GET\",\"url\":\"http://foo.com/health\",\"status\":500,\"size\":45,\"duration\":2,\"message\":\"GET /health\"}\n", string(actual))
}
