// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

// SourceProvider interface
type SourceProvider interface {
	Get(resource *Resource) (*Resource, error)
	Set(resource *Resource) error
	Del(resource *Resource) bool
}
