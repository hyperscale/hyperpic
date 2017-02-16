package main

import (
	"net/http"

	"github.com/euskadi31/image-service/memfs"
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
