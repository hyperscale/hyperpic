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
