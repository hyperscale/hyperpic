// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package provider

import (
	"github.com/hyperscale/hyperpic/image"
)

// CacheProvider interface
//go:generate mockery -case=underscore -inpkg -name=CacheProvider
type CacheProvider interface {
	Get(resource *image.Resource) (*image.Resource, error)

	Set(resource *image.Resource) error

	Del(resource *image.Resource) error
}
