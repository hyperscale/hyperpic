// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

// ImageSupportConfiguration struct
type ImageSupportConfiguration struct {
	Extensions map[string]interface{}
}

// IsExtSupported return true if ext is supported
func (c ImageSupportConfiguration) IsExtSupported(ext string) bool {
	enable, ok := c.Extensions[ext]

	return (ok && enable.(bool))
}
