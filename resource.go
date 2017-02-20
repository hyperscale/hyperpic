// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import "time"

// Resource struct
type Resource struct {
	Name       string
	Path       string
	Options    *ImageOptions
	MimeType   string
	ModifiedAt time.Time
	Body       []byte
}
