// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package fsutil

import (
	"strings"
)

// ContainsDotDot return true if string contains ..
func ContainsDotDot(v string) bool {
	if !strings.Contains(v, "..") {
		return false
	}

	for _, ent := range strings.FieldsFunc(v, isSlashRune) {
		if ent == ".." {
			return true
		}
	}

	return false
}

func isSlashRune(r rune) bool {
	return r == '/' || r == '\\'
}
