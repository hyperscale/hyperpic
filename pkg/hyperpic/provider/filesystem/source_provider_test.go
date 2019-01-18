// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package filesystem

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/hyperscale/hyperpic/pkg/hyperpic/image"
	"github.com/stretchr/testify/assert"
)

func TestSourceProvider(t *testing.T) {
	dir, err := ioutil.TempDir("", "source-provider-test")
	assert.NoError(t, err)

	defer os.RemoveAll(dir)

	p := NewSourceProvider(&SourceConfiguration{
		Path: dir,
	})

	err = p.Del(&image.Resource{
		Path: "/../../",
	})
	assert.Error(t, err)

	err = p.Del(&image.Resource{
		Path: "/test.jpg",
	})
	assert.Error(t, err)

	res, err := p.Get(&image.Resource{
		Path: "/test.jpg",
	})
	assert.Error(t, err)
	assert.Nil(t, res)

	res, err = p.Get(&image.Resource{
		Path: "/",
	})
	assert.Error(t, err)
	assert.Nil(t, res)

	res, err = p.Get(&image.Resource{
		Path: "/../../",
	})
	assert.Error(t, err)
	assert.Nil(t, res)

	err = p.Set(&image.Resource{
		Path: "/test/",
	})
	assert.Error(t, err)

	body, err := ioutil.ReadFile("../../../../_resources/demo/kayaks.jpg")
	assert.NoError(t, err)

	err = p.Set(&image.Resource{
		Path: "/test.jpg",
		Body: body,
		Size: len(body),
	})
	assert.NoError(t, err)

	res, err = p.Get(&image.Resource{
		Path: "/test.jpg",
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, body, res.Body)
	assert.Equal(t, "", res.MimeType)
	assert.Equal(t, "test.jpg", res.Name)
	assert.Equal(t, "/test.jpg", res.Path)
}
