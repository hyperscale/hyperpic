// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package filesystem

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/hyperscale/hyperpic/pkg/hyperpic/image"
	"github.com/stretchr/testify/assert"
)

func TestCacheProvider(t *testing.T) {
	dir, err := ioutil.TempDir("", "cache-provider-test")
	assert.NoError(t, err)

	defer os.RemoveAll(dir)

	p := NewCacheProvider(&CacheConfiguration{
		Path:          dir,
		LifeTime:      1 * time.Second,
		CleanInterval: 2 * time.Second,
	})

	body, err := os.ReadFile("../../../../_resources/demo/kayaks.jpg")
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

	res, err := p.Get(&image.Resource{
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

	err = p.Del(&image.Resource{
		Path: "/..",
	})
	assert.Error(t, err)

	err = p.Del(&image.Resource{
		Path: "/kayaks.jpg",
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
