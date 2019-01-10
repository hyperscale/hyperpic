// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controllers

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	server "github.com/euskadi31/go-server"
	"github.com/hyperscale/hyperpic/config"
	"github.com/hyperscale/hyperpic/image"
	"github.com/hyperscale/hyperpic/provider"
	"github.com/hyperscale/hyperpic/provider/filesystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestImageControllerGetNotFoundImage(t *testing.T) {
	cfg := &config.Configuration{
		Image: &config.ImageConfiguration{
			Source: &config.ImageSourceConfiguration{
				FS: &config.SourceFSConfiguration{
					Path: "../_resources/demo",
				},
			},
			Support: &config.ImageSupportConfiguration{
				Extensions: map[string]interface{}{
					"jpg":  true,
					"jpeg": true,
					"png":  true,
					"webp": true,
				},
			},
		},
	}

	optionsParser := image.NewOptionParser()
	sourceProvider := filesystem.NewSourceProvider(cfg.Image.Source.FS)
	cacheProvider := &provider.MockCacheProvider{}

	cacheProvider.On("Get", mock.MatchedBy(func(res *image.Resource) bool {
		if res.Path != "/not-found.jpg" {
			return false
		}

		return true
	})).Return(nil, errors.New("not exist"))

	controller, _ := NewImageController(cfg, optionsParser, sourceProvider, cacheProvider)

	router := server.NewRouter()

	router.AddController(controller)

	req := httptest.NewRequest(http.MethodGet, "/not-found.jpg?w=40&h=40&q=85&fm=webp", nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestImageControllerGetImageWithSourceProviderError(t *testing.T) {
	cfg := &config.Configuration{
		Image: &config.ImageConfiguration{
			Support: &config.ImageSupportConfiguration{
				Extensions: map[string]interface{}{
					"jpg":  true,
					"jpeg": true,
					"png":  true,
					"webp": true,
				},
			},
		},
	}

	optionsParser := image.NewOptionParser()

	sourceProvider := &provider.MockSourceProvider{}

	sourceProvider.On("Get", mock.MatchedBy(func(res *image.Resource) bool {
		if res.Path != "/not-found.jpg" {
			return false
		}

		return true
	})).Return(nil, errors.New("fail"))

	cacheProvider := &provider.MockCacheProvider{}

	cacheProvider.On("Get", mock.MatchedBy(func(res *image.Resource) bool {
		if res.Path != "/not-found.jpg" {
			return false
		}

		return true
	})).Return(nil, errors.New("not exist"))

	controller, _ := NewImageController(cfg, optionsParser, sourceProvider, cacheProvider)

	router := server.NewRouter()

	router.AddController(controller)

	req := httptest.NewRequest(http.MethodGet, "/not-found.jpg?w=40&h=40&q=85&fm=webp", nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestImageControllerGetImageWithProccessImageError(t *testing.T) {
	cfg := &config.Configuration{
		Image: &config.ImageConfiguration{
			Support: &config.ImageSupportConfiguration{
				Extensions: map[string]interface{}{
					"jpg":  true,
					"jpeg": true,
					"png":  true,
					"webp": true,
				},
			},
		},
	}

	optionsParser := image.NewOptionParser()

	sourceProvider := &provider.MockSourceProvider{}

	sourceProvider.On("Get", mock.MatchedBy(func(res *image.Resource) bool {
		if res.Path != "/kayaks.jpg" {
			return false
		}

		return true
	})).Return(&image.Resource{
		Path:       "/kayaks.jpg",
		Body:       nil,
		ModifiedAt: time.Now(),
		Options:    nil,
	}, nil)

	cacheProvider := &provider.MockCacheProvider{}

	cacheProvider.On("Get", mock.MatchedBy(func(res *image.Resource) bool {
		if res.Path != "/kayaks.jpg" {
			return false
		}

		return true
	})).Return(nil, errors.New("not exist"))

	controller, _ := NewImageController(cfg, optionsParser, sourceProvider, cacheProvider)

	router := server.NewRouter()

	router.AddController(controller)

	req := httptest.NewRequest(http.MethodGet, "/kayaks.jpg?w=40&h=40&q=85&fm=webp", nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestImageControllerGetImageWithFormatInQueryString(t *testing.T) {
	cfg := &config.Configuration{
		Image: &config.ImageConfiguration{
			Source: &config.ImageSourceConfiguration{
				FS: &config.SourceFSConfiguration{
					Path: "../_resources/demo",
				},
			},
			Support: &config.ImageSupportConfiguration{
				Extensions: map[string]interface{}{
					"jpg":  true,
					"jpeg": true,
					"png":  true,
					"webp": true,
				},
			},
		},
	}

	optionsParser := image.NewOptionParser()
	sourceProvider := filesystem.NewSourceProvider(cfg.Image.Source.FS)
	cacheProvider := &provider.MockCacheProvider{}

	cacheProvider.On("Get", mock.MatchedBy(func(res *image.Resource) bool {
		if res.Path != "/kayaks.jpg" {
			return false
		}

		return true
	})).Return(nil, errors.New("not exist"))

	cacheProvider.On("Set", mock.MatchedBy(func(res *image.Resource) bool {
		if res.Path != "/kayaks.jpg" {
			return false
		}

		if res.Name != "kayaks.jpg" {
			return false
		}

		if res.MimeType != "image/webp" {
			return false
		}

		return true
	})).Return(errors.New("fail"))

	controller, _ := NewImageController(cfg, optionsParser, sourceProvider, cacheProvider)

	router := server.NewRouter()

	router.AddController(controller)

	req := httptest.NewRequest(http.MethodGet, "/kayaks.jpg?w=40&h=40&q=85&fm=webp", nil)
	req.Header.Set("Accept", "*/*")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "image/webp", resp.Header.Get("Content-Type"))
	assert.Equal(t, "source", resp.Header.Get("X-Image-From"))
}

func TestImageControllerGetImageInCache(t *testing.T) {
	cfg := &config.Configuration{
		Image: &config.ImageConfiguration{
			Source: &config.ImageSourceConfiguration{
				FS: &config.SourceFSConfiguration{
					Path: "../_resources/demo",
				},
			},
			Support: &config.ImageSupportConfiguration{
				Extensions: map[string]interface{}{
					"jpg":  true,
					"jpeg": true,
					"png":  true,
					"webp": true,
				},
			},
		},
	}

	optionsParser := image.NewOptionParser()
	sourceProvider := filesystem.NewSourceProvider(cfg.Image.Source.FS)
	cacheProvider := &provider.MockCacheProvider{}

	data, err := ioutil.ReadFile("../_resources/demo/kayaks.jpg")
	assert.NoError(t, err)

	var opts *image.Options

	cacheProvider.On("Get", mock.MatchedBy(func(res *image.Resource) bool {
		if res.Path != "/kayaks.jpg" {
			return false
		}

		return true
	})).Return(&image.Resource{
		Path:       "/kayaks.jpg",
		Body:       data,
		ModifiedAt: time.Now(),
		Options:    opts,
	}, nil).Run(func(args mock.Arguments) {
		res := args.Get(0).(*image.Resource)

		opts = res.Options
	})

	controller, _ := NewImageController(cfg, optionsParser, sourceProvider, cacheProvider)

	router := server.NewRouter()

	router.AddController(controller)

	req := httptest.NewRequest(http.MethodGet, "/kayaks.jpg?w=40&h=40&q=85&fm=webp", nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "image/webp", resp.Header.Get("Content-Type"))
	assert.Equal(t, "cache", resp.Header.Get("X-Image-From"))
}

func TestImageControllerGetHandlerWithoutOptionsParserInContext(t *testing.T) {
	controller := &imageController{}

	req := httptest.NewRequest(http.MethodGet, "/kayaks.jpg?w=40&h=40&q=85&fm=webp", nil)

	w := httptest.NewRecorder()

	controller.getHandler(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestImageControllerDeleteImageCache(t *testing.T) {
	cfg := &config.Configuration{
		Auth: &config.AuthConfiguration{
			Secret: "foo",
		},
		Image: &config.ImageConfiguration{
			Support: &config.ImageSupportConfiguration{
				Extensions: map[string]interface{}{
					"jpg":  true,
					"jpeg": true,
					"png":  true,
					"webp": true,
				},
			},
		},
	}

	optionsParser := image.NewOptionParser()
	cacheProvider := &provider.MockCacheProvider{}

	cacheProvider.On("Del", mock.MatchedBy(func(res *image.Resource) bool {
		if res.Path != "/kayaks.jpg" {
			return false
		}

		return true
	})).Return(nil)

	controller, _ := NewImageController(cfg, optionsParser, nil, cacheProvider)

	router := server.NewRouter()

	router.AddController(controller)

	req := httptest.NewRequest(http.MethodDelete, "/kayaks.jpg?from=cache", nil)
	req.Header.Set("Authorization", "Bearer foo")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	actuel, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.JSONEq(t, `{"cache": true, "source": false}`, string(actuel))
}

func TestImageControllerDeleteImageSource(t *testing.T) {
	cfg := &config.Configuration{
		Auth: &config.AuthConfiguration{
			Secret: "foo",
		},
		Image: &config.ImageConfiguration{
			Support: &config.ImageSupportConfiguration{
				Extensions: map[string]interface{}{
					"jpg":  true,
					"jpeg": true,
					"png":  true,
					"webp": true,
				},
			},
		},
	}

	optionsParser := image.NewOptionParser()

	sourceProvider := &provider.MockSourceProvider{}

	sourceProvider.On("Del", mock.MatchedBy(func(res *image.Resource) bool {
		if res.Path != "/kayaks.jpg" {
			return false
		}

		return true
	})).Return(nil)

	cacheProvider := &provider.MockCacheProvider{}

	cacheProvider.On("Del", mock.MatchedBy(func(res *image.Resource) bool {
		if res.Path != "/kayaks.jpg" {
			return false
		}

		return true
	})).Return(nil)

	controller, _ := NewImageController(cfg, optionsParser, sourceProvider, cacheProvider)

	router := server.NewRouter()

	router.AddController(controller)

	req := httptest.NewRequest(http.MethodDelete, "/kayaks.jpg?from=source", nil)
	req.Header.Set("Authorization", "Bearer foo")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	actuel, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.JSONEq(t, `{"cache": true, "source": true}`, string(actuel))
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestImageControllerPostImageWithBadBody(t *testing.T) {
	cfg := &config.Configuration{
		Auth: &config.AuthConfiguration{
			Secret: "foo",
		},
		Image: &config.ImageConfiguration{
			Support: &config.ImageSupportConfiguration{
				Extensions: map[string]interface{}{
					"jpg":  true,
					"jpeg": true,
					"png":  true,
					"webp": true,
				},
			},
		},
	}

	optionsParser := image.NewOptionParser()

	controller, _ := NewImageController(cfg, optionsParser, nil, nil)

	router := server.NewRouter()

	router.AddController(controller)

	req := httptest.NewRequest(http.MethodPost, "/kayaks.jpg", errReader(0))
	req.Header.Set("Authorization", "Bearer foo")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestImageControllerParseImageFileFromRequestWithEmptyBody(t *testing.T) {
	controller := &imageController{}

	req := httptest.NewRequest(http.MethodPost, "/kayaks.jpg", nil)
	req.Body = nil

	data, err := controller.parseImageFileFromRequest(req)
	assert.Nil(t, data)
	assert.Error(t, err)
}

func TestImageControllerParseImageFileFromRequestWithBadContentType(t *testing.T) {
	controller := &imageController{}

	req := httptest.NewRequest(http.MethodPost, "/kayaks.jpg", nil)
	req.Header.Set("Content-Type", "E34535;;/")

	data, err := controller.parseImageFileFromRequest(req)
	assert.Nil(t, data)
	assert.Error(t, err)
}

func TestImageControllerParseImageFileFromRequestMultipartWithBadField(t *testing.T) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	src, err := ioutil.ReadFile("../_resources/demo/kayaks.jpg")
	assert.NoError(t, err)

	fw, err := w.CreateFormFile("bad", "../_resources/demo/kayaks.jpg")
	assert.NoError(t, err)

	_, err = io.Copy(fw, bytes.NewReader(src))
	assert.NoError(t, err)

	w.Close()

	controller := &imageController{}

	req := httptest.NewRequest(http.MethodPost, "/kayaks.jpg", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())

	data, err := controller.parseImageFileFromRequest(req)
	assert.Nil(t, data)
	assert.Error(t, err)
}

func TestImageControllerParseImageFileFromRequestMultipart(t *testing.T) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	src, err := ioutil.ReadFile("../_resources/demo/kayaks.jpg")
	assert.NoError(t, err)

	fw, err := w.CreateFormFile("image", "../_resources/demo/kayaks.jpg")
	assert.NoError(t, err)

	_, err = io.Copy(fw, bytes.NewReader(src))
	assert.NoError(t, err)

	w.Close()

	controller := &imageController{}

	req := httptest.NewRequest(http.MethodPost, "/kayaks.jpg", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())

	data, err := controller.parseImageFileFromRequest(req)
	assert.Equal(t, len(src), len(data))
	assert.NoError(t, err)
}

func TestImageControllerPostImageWithFailOnSourceProviderSet(t *testing.T) {
	data, err := ioutil.ReadFile("../_resources/demo/kayaks.jpg")
	assert.NoError(t, err)

	cfg := &config.Configuration{
		Auth: &config.AuthConfiguration{
			Secret: "foo",
		},
		Image: &config.ImageConfiguration{
			Support: &config.ImageSupportConfiguration{
				Extensions: map[string]interface{}{
					"jpg":  true,
					"jpeg": true,
					"png":  true,
					"webp": true,
				},
			},
		},
	}

	optionsParser := image.NewOptionParser()

	sourceProvider := &provider.MockSourceProvider{}

	sourceProvider.On("Set", mock.MatchedBy(func(res *image.Resource) bool {
		if res.Path != "/kayaks.jpg" {
			return false
		}

		if bytes.Compare(res.Body, data) != 0 {
			return false
		}

		return true
	})).Return(errors.New("fail"))

	cacheProvider := &provider.MockCacheProvider{}

	controller, _ := NewImageController(cfg, optionsParser, sourceProvider, cacheProvider)

	router := server.NewRouter()

	router.AddController(controller)

	req := httptest.NewRequest(http.MethodPost, "/kayaks.jpg", bytes.NewReader(data))
	req.Header.Set("Authorization", "Bearer foo")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestImageControllerPostImage(t *testing.T) {
	data, err := ioutil.ReadFile("../_resources/demo/kayaks.jpg")
	assert.NoError(t, err)

	cfg := &config.Configuration{
		Auth: &config.AuthConfiguration{
			Secret: "foo",
		},
		Image: &config.ImageConfiguration{
			Support: &config.ImageSupportConfiguration{
				Extensions: map[string]interface{}{
					"jpg":  true,
					"jpeg": true,
					"png":  true,
					"webp": true,
				},
			},
		},
	}

	optionsParser := image.NewOptionParser()

	sourceProvider := &provider.MockSourceProvider{}

	sourceProvider.On("Set", mock.MatchedBy(func(res *image.Resource) bool {
		if res.Path != "/kayaks.jpg" {
			return false
		}

		if bytes.Compare(res.Body, data) != 0 {
			return false
		}

		return true
	})).Return(nil)

	cacheProvider := &provider.MockCacheProvider{}

	cacheProvider.On("Del", mock.MatchedBy(func(res *image.Resource) bool {
		if res.Path != "/kayaks.jpg" {
			return false
		}

		return true
	})).Return(errors.New("fail")).Maybe()

	controller, _ := NewImageController(cfg, optionsParser, sourceProvider, cacheProvider)

	router := server.NewRouter()

	router.AddController(controller)

	req := httptest.NewRequest(http.MethodPost, "/kayaks.jpg", bytes.NewReader(data))
	req.Header.Set("Authorization", "Bearer foo")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	actuel, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.JSONEq(t, `{"file": "/kayaks.jpg", "size": 256355, "type": "image/jpeg", "hash": "7f267455b649191f5e4f4d56b9766c6c"}`, string(actuel))

	cacheProvider.AssertExpectations(t)
	sourceProvider.AssertExpectations(t)
}
