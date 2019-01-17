// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controller

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	server "github.com/euskadi31/go-server"
	"github.com/hyperscale/hyperpic/cmd/hyperpic/app/config"
	"github.com/hyperscale/hyperpic/pkg/hyperpic/image"
	"github.com/hyperscale/hyperpic/pkg/hyperpic/provider"
	"github.com/hyperscale/hyperpic/pkg/hyperpic/provider/filesystem"
	"github.com/hyperscale/hyperpic/pkg/hyperpic/provider/memory"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestImageControllerGetNotFoundImage(t *testing.T) {
	cfg := &config.Configuration{
		Image: &config.ImageConfiguration{
			Source: &config.ImageSourceConfiguration{
				FS: &filesystem.SourceConfiguration{
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

	imageProcessor := &image.MockProcessor{}

	controller := NewImageController(cfg, optionsParser, imageProcessor, sourceProvider, cacheProvider)

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

	imageProcessor := &image.MockProcessor{}

	controller := NewImageController(cfg, optionsParser, imageProcessor, sourceProvider, cacheProvider)

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

	imageProcessor := &image.MockProcessor{}

	imageProcessor.On("ProcessImage", mock.MatchedBy(func(res *image.Resource) bool {
		if res.Path != "/kayaks.jpg" {
			return false
		}

		return true
	})).Return(errors.New("foo"))

	controller := NewImageController(cfg, optionsParser, imageProcessor, sourceProvider, cacheProvider)

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
				FS: &filesystem.SourceConfiguration{
					Path: "../../../../_resources/demo",
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

	imageProcessor := image.NewProcessor()

	controller := NewImageController(cfg, optionsParser, imageProcessor, sourceProvider, cacheProvider)

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
				FS: &filesystem.SourceConfiguration{
					Path: "../../../../_resources/demo",
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

	data, err := ioutil.ReadFile("../../../../_resources/demo/kayaks.jpg")
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

	imageProcessor := &image.MockProcessor{}

	controller := NewImageController(cfg, optionsParser, imageProcessor, sourceProvider, cacheProvider)

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

	imageProcessor := &image.MockProcessor{}

	controller := NewImageController(cfg, optionsParser, imageProcessor, nil, cacheProvider)

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

	imageProcessor := &image.MockProcessor{}

	controller := NewImageController(cfg, optionsParser, imageProcessor, sourceProvider, cacheProvider)

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
			Source: &config.ImageSourceConfiguration{
				MaxSize: 10 << 20,
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

	imageProcessor := &image.MockProcessor{}

	controller := NewImageController(cfg, optionsParser, imageProcessor, nil, nil)

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

	w := httptest.NewRecorder()

	data, err := controller.parseImageFileFromRequest(w, req)
	assert.Nil(t, data)
	assert.Error(t, err)
}

func TestImageControllerParseImageFileFromRequestWithBadContentType(t *testing.T) {
	controller := &imageController{
		cfg: &config.Configuration{
			Image: &config.ImageConfiguration{
				Source: &config.ImageSourceConfiguration{
					MaxSize: 10 << 20,
				},
			},
		},
	}

	req := httptest.NewRequest(http.MethodPost, "/kayaks.jpg", nil)
	req.Header.Set("Content-Type", "E34535;;/")

	w := httptest.NewRecorder()

	data, err := controller.parseImageFileFromRequest(w, req)
	assert.Nil(t, data)
	assert.Error(t, err)
}

func TestImageControllerParseImageFileFromRequestMultipartWithBadField(t *testing.T) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	src, err := ioutil.ReadFile("../../../../_resources/demo/kayaks.jpg")
	assert.NoError(t, err)

	fw, err := w.CreateFormFile("bad", "../../../../_resources/demo/kayaks.jpg")
	assert.NoError(t, err)

	_, err = io.Copy(fw, bytes.NewReader(src))
	assert.NoError(t, err)

	w.Close()

	controller := &imageController{
		cfg: &config.Configuration{
			Image: &config.ImageConfiguration{
				Source: &config.ImageSourceConfiguration{
					MaxSize: 10 << 20,
				},
			},
		},
	}

	req := httptest.NewRequest(http.MethodPost, "/kayaks.jpg", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())

	wr := httptest.NewRecorder()

	data, err := controller.parseImageFileFromRequest(wr, req)
	assert.Nil(t, data)
	assert.Error(t, err)
}

func TestImageControllerParseImageFileFromRequestMultipart(t *testing.T) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	src, err := ioutil.ReadFile("../../../../_resources/demo/kayaks.jpg")
	assert.NoError(t, err)

	fw, err := w.CreateFormFile("image", "../../../../_resources/demo/kayaks.jpg")
	assert.NoError(t, err)

	_, err = io.Copy(fw, bytes.NewReader(src))
	assert.NoError(t, err)

	w.Close()

	controller := &imageController{
		cfg: &config.Configuration{
			Image: &config.ImageConfiguration{
				Source: &config.ImageSourceConfiguration{
					MaxSize: 10 << 20,
				},
			},
		},
	}

	req := httptest.NewRequest(http.MethodPost, "/kayaks.jpg", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())

	wr := httptest.NewRecorder()

	data, err := controller.parseImageFileFromRequest(wr, req)
	assert.Equal(t, len(src), len(data))
	assert.NoError(t, err)
}

func TestImageControllerParseImageFileFromRequestMultipartWithTooBigFile(t *testing.T) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	src, err := ioutil.ReadFile("../../../../_resources/demo/kayaks.jpg")
	assert.NoError(t, err)

	fw, err := w.CreateFormFile("image", "../../../../_resources/demo/kayaks.jpg")
	assert.NoError(t, err)

	_, err = io.Copy(fw, bytes.NewReader(src))
	assert.NoError(t, err)

	w.Close()

	controller := &imageController{
		cfg: &config.Configuration{
			Image: &config.ImageConfiguration{
				Source: &config.ImageSourceConfiguration{
					MaxSize: 100,
				},
			},
		},
	}

	req := httptest.NewRequest(http.MethodPost, "/kayaks.jpg", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())

	wr := httptest.NewRecorder()

	data, err := controller.parseImageFileFromRequest(wr, req)
	assert.Nil(t, data)
	assert.EqualError(t, err, "http: request body too large")
}

func TestImageControllerPostImageWithFailOnSourceProviderSet(t *testing.T) {
	data, err := ioutil.ReadFile("../../../../_resources/demo/kayaks.jpg")
	assert.NoError(t, err)

	cfg := &config.Configuration{
		Auth: &config.AuthConfiguration{
			Secret: "foo",
		},
		Image: &config.ImageConfiguration{
			Source: &config.ImageSourceConfiguration{
				MaxSize: 10 << 20,
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

	imageProcessor := &image.MockProcessor{}

	controller := NewImageController(cfg, optionsParser, imageProcessor, sourceProvider, cacheProvider)

	router := server.NewRouter()

	router.AddController(controller)

	req := httptest.NewRequest(http.MethodPost, "/kayaks.jpg", bytes.NewReader(data))
	req.Header.Set("Authorization", "Bearer foo")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestImageControllerPostImageWithTooLargeError(t *testing.T) {
	data, err := ioutil.ReadFile("../../../../_resources/demo/kayaks.jpg")
	assert.NoError(t, err)

	cfg := &config.Configuration{
		Auth: &config.AuthConfiguration{
			Secret: "foo",
		},
		Image: &config.ImageConfiguration{
			Source: &config.ImageSourceConfiguration{
				MaxSize: 1 << 10,
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

	imageProcessor := &image.MockProcessor{}

	controller := NewImageController(cfg, optionsParser, imageProcessor, nil, nil)

	router := server.NewRouter()

	router.AddController(controller)

	req := httptest.NewRequest(http.MethodPost, "/kayaks.jpg", bytes.NewReader(data))
	req.Header.Set("Authorization", "Bearer foo")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusRequestEntityTooLarge, resp.StatusCode)
}

func TestImageControllerPostImage(t *testing.T) {
	data, err := ioutil.ReadFile("../../../../_resources/demo/kayaks.jpg")
	assert.NoError(t, err)

	cfg := &config.Configuration{
		Auth: &config.AuthConfiguration{
			Secret: "foo",
		},
		Image: &config.ImageConfiguration{
			Source: &config.ImageSourceConfiguration{
				MaxSize: 10 << 20,
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

	imageProcessor := &image.MockProcessor{}

	controller := NewImageController(cfg, optionsParser, imageProcessor, sourceProvider, cacheProvider)

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
	assert.JSONEq(t, `{"file": "/kayaks.jpg", "size": 256355, "type": "image/jpeg", "hash": "8138fdd61f7d8b3ac0d0f11cd2fe994fe37f8657cb93f6e8f818606294c7079e"}`, string(actuel))

	cacheProvider.AssertExpectations(t)
	sourceProvider.AssertExpectations(t)
}

func BenchmarkProcessImageNoCache(b *testing.B) {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	time.Sleep(100 * time.Millisecond)

	cfg := &config.Configuration{
		Image: &config.ImageConfiguration{
			Source: &config.ImageSourceConfiguration{
				FS: &filesystem.SourceConfiguration{
					Path: "../../../../_resources/demo",
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

	imageProcessor := image.NewProcessor()

	controller := NewImageController(cfg, optionsParser, imageProcessor, sourceProvider, cacheProvider)

	router := server.NewRouter()

	router.AddController(controller)

	req := httptest.NewRequest(http.MethodGet, "/kayaks.jpg?w=40&h=40&q=85&fm=webp", nil)
	req.Header.Set("Accept", "*/*")

	for n := 0; n < b.N; n++ {
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
	}

	time.Sleep(100 * time.Millisecond)

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func BenchmarkProcessImageWithFSCache(b *testing.B) {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	time.Sleep(100 * time.Millisecond)

	dir, err := ioutil.TempDir("", "cache-provider-bench")
	assert.NoError(b, err)

	defer os.RemoveAll(dir)

	cfg := &config.Configuration{
		Image: &config.ImageConfiguration{
			Source: &config.ImageSourceConfiguration{
				FS: &filesystem.SourceConfiguration{
					Path: "../../../../_resources/demo",
				},
			},
			Cache: &config.ImageCacheConfiguration{
				FS: &filesystem.CacheConfiguration{
					Path:          dir,
					LifeTime:      1 * time.Minute,
					CleanInterval: 1 * time.Minute,
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
	cacheProvider := filesystem.NewCacheProvider(cfg.Image.Cache.FS)

	imageProcessor := image.NewProcessor()

	controller := NewImageController(cfg, optionsParser, imageProcessor, sourceProvider, cacheProvider)

	router := server.NewRouter()

	router.AddController(controller)

	req := httptest.NewRequest(http.MethodGet, "/kayaks.jpg?w=40&h=40&q=85&fm=webp", nil)
	req.Header.Set("Accept", "*/*")

	for n := 0; n < b.N; n++ {
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
	}

	time.Sleep(100 * time.Millisecond)

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func BenchmarkProcessImageWithMemoryCache(b *testing.B) {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	time.Sleep(100 * time.Millisecond)

	cfg := &config.Configuration{
		Image: &config.ImageConfiguration{
			Source: &config.ImageSourceConfiguration{
				FS: &filesystem.SourceConfiguration{
					Path: "../../../../_resources/demo",
				},
			},
			Cache: &config.ImageCacheConfiguration{
				Memory: &memory.CacheConfiguration{
					LifeTime:      1 * time.Minute,
					CleanInterval: 1 * time.Minute,
					MemoryLimit:   10 << 20,
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
	cacheProvider := memory.NewCacheProvider(cfg.Image.Cache.Memory)

	imageProcessor := image.NewProcessor()

	controller := NewImageController(cfg, optionsParser, imageProcessor, sourceProvider, cacheProvider)

	router := server.NewRouter()

	router.AddController(controller)

	req := httptest.NewRequest(http.MethodGet, "/kayaks.jpg?w=40&h=40&q=85&fm=webp", nil)
	req.Header.Set("Accept", "*/*")

	for n := 0; n < b.N; n++ {
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
	}

	time.Sleep(100 * time.Millisecond)

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}
