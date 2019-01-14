// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package httputil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorMessage(t *testing.T) {
	err := ErrorMessage{
		Code:    100,
		Message: "Test",
	}

	assert.Equal(t, 100, err.GetCode())
	assert.Equal(t, "Test", err.GetMessage())
	assert.Equal(t, "Test", err.Error())
}
