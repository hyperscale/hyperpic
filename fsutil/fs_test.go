// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package fsutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsDotDot(t *testing.T) {
	expects := []struct {
		Value    string
		Expected bool
	}{
		{
			Value:    "../../",
			Expected: true,
		},
		{
			Value:    "/my/image/../../",
			Expected: true,
		},
		{
			Value:    "/my/image.jpg",
			Expected: false,
		},
		{
			Value:    "/my/ima..ge.jpg",
			Expected: false,
		},
	}

	for _, item := range expects {
		if item.Expected {
			assert.True(t, ContainsDotDot(item.Value))
		} else {
			assert.False(t, ContainsDotDot(item.Value))
		}
	}
}
