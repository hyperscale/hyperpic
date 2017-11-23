// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package httputil

import (
	"net/http"

	"github.com/hyperscale/hyperpic/image"
	"github.com/hyperscale/hyperpic/memfs"
)

// ServeImage from resource
func ServeImage(w http.ResponseWriter, r *http.Request, resource *image.Resource) {
	http.ServeContent(
		w,
		r,
		resource.Name,
		resource.ModifiedAt,
		memfs.NewBuffer(&resource.Body),
	)
}
