// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"net/http"

	"github.com/hyperscale/hyperpic/memfs"
)

func ServeImage(w http.ResponseWriter, r *http.Request, resource *Resource) {
	http.ServeContent(
		w,
		r,
		resource.Name,
		resource.ModifiedAt,
		memfs.NewBuffer(&resource.Body),
	)
}
