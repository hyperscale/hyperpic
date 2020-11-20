// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package memory

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/hyperscale/hyperpic/pkg/hyperpic/image"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestCacheProvider(t *testing.T) {
	p := NewCacheProvider(&CacheConfiguration{
		LifeTime:      1 * time.Second,
		CleanInterval: 2 * time.Second,
		MemoryLimit:   256356,
	})

	body, err := ioutil.ReadFile("../../../../_resources/demo/kayaks.jpg")
	assert.NoError(t, err)

	err = p.Set(&image.Resource{
		Path: "/kayaks.jpg",
		Body: body,
		Size: len(body),
		Options: &image.Options{
			Width:  200,
			Height: 200,
		},
	})
	assert.NoError(t, err)

	err = p.Set(&image.Resource{
		Path: "/kayaks-2.jpg",
		Body: body,
		Size: len(body),
		Options: &image.Options{
			Width:  200,
			Height: 200,
		},
	})
	assert.EqualError(t, err, "memory cache provider: allowed memory size of 256356 bytes exhausted")

	res, err := p.Get(&image.Resource{
		Path: "/kayaks-2.jpg",
		Options: &image.Options{
			Width:  200,
			Height: 200,
		},
	})
	assert.EqualError(t, err, "file does not exist in cache")
	assert.Nil(t, res)

	res, err = p.Get(&image.Resource{
		Path: "/..",
		Options: &image.Options{
			Width:  100,
			Height: 100,
		},
	})
	assert.Error(t, err)
	assert.Nil(t, res)

	res, err = p.Get(&image.Resource{
		Path: "/kayaks.jpg",
		Options: &image.Options{
			Width:  100,
			Height: 100,
		},
	})
	assert.Error(t, err)
	assert.Nil(t, res)

	res, err = p.Get(&image.Resource{
		Path: "/kayaks.jpg",
		Options: &image.Options{
			Width:  200,
			Height: 200,
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, body, res.Body)
	assert.Equal(t, len(body), res.Size)

	err = p.Del(&image.Resource{
		Path: "/..",
	})
	assert.Error(t, err)

	err = p.Del(&image.Resource{
		Path: "/kayaks.jpg",
	})
	assert.NoError(t, err)

	err = p.Del(&image.Resource{
		Path: "/kayaks-3.jpg",
	})
	assert.NoError(t, err)

	err = p.Set(&image.Resource{
		Path: "/kayaks.jpg",
		Body: body,
		Size: len(body),
		Options: &image.Options{
			Width:  200,
			Height: 200,
		},
	})
	assert.NoError(t, err)

	time.Sleep(3 * time.Second)

	res, err = p.Get(&image.Resource{
		Path: "/kayaks.jpg",
		Options: &image.Options{
			Width:  200,
			Height: 200,
		},
	})
	assert.Error(t, err)
	assert.Nil(t, res)

	time.Sleep(100 * time.Millisecond)
}

func BenchmarkCacheProviderSet(b *testing.B) {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	time.Sleep(100 * time.Millisecond)

	p := NewCacheProvider(&CacheConfiguration{
		LifeTime:      30 * time.Second,
		CleanInterval: 35 * time.Second,
		MemoryLimit:   10 << 20,
	})

	body, err := ioutil.ReadFile("../../../../_resources/demo/kayaks.jpg")
	assert.NoError(b, err)

	for n := 0; n < b.N; n++ {
		err := p.Set(&image.Resource{
			Path: fmt.Sprintf("/kayaks-%d.jpg", n),
			Body: body,
			Size: 0, // hack for benchmark
			Options: &image.Options{
				Width:  200,
				Height: 200,
			},
		})
		assert.NoError(b, err)
	}

	time.Sleep(100 * time.Millisecond)

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func BenchmarkCacheProviderGet(b *testing.B) {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	time.Sleep(100 * time.Millisecond)

	p := NewCacheProvider(&CacheConfiguration{
		LifeTime:      30 * time.Second,
		CleanInterval: 35 * time.Second,
		MemoryLimit:   10 << 20,
	})

	body, err := ioutil.ReadFile("../../../../_resources/demo/kayaks.jpg")
	assert.NoError(b, err)

	err = p.Set(&image.Resource{
		Path: "/kayaks.jpg",
		Body: body,
		Size: len(body),
		Options: &image.Options{
			Width:  200,
			Height: 200,
		},
	})
	assert.NoError(b, err)

	for n := 0; n < b.N; n++ {
		_, err = p.Get(&image.Resource{
			Path: "/kayaks.jpg",
			Options: &image.Options{
				Width:  200,
				Height: 200,
			},
		})
		assert.NoError(b, err)
	}

	time.Sleep(100 * time.Millisecond)

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}
