// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package middlewares

import "errors"

var (
	errContextIsNull     = errors.New("the context is null")
	errNotFountInContext = errors.New("the entry is not found in context")
)

type key int

const (
	optionsKey key = iota
)
