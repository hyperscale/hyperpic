// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

// CacheProvider interface
type CacheProvider interface {
	Get(resource *Resource) (*Resource, bool)

	Set(resource *Resource) error

	Del(resource *Resource) bool
}
