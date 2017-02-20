// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"strings"

	bimg "gopkg.in/h2non/bimg.v1"
)

// ExtractImageTypeFromMime returns the MIME image type.
func ExtractImageTypeFromMime(mime string) string {
	mime = strings.Split(mime, ";")[0]
	parts := strings.Split(mime, "/")

	if len(parts) < 2 {
		return ""
	}

	name := strings.Split(parts[1], "+")[0]

	return strings.ToLower(name)
}

// IsImageMimeTypeSupported returns true if the image MIME
// type is supported by bimg.
func IsImageMimeTypeSupported(mime string) bool {
	format := ExtractImageTypeFromMime(mime)

	// Some payloads may expose the MIME type for SVG as text/xml
	if format == "xml" {
		format = "svg"
	}

	return bimg.IsTypeNameSupported(format)
}

// ImageType returns the image type based on the given image type alias.
func ImageType(name string) bimg.ImageType {
	ext := strings.ToLower(name)

	switch ext {
	case "jpeg":
		return bimg.JPEG
	case "jpg":
		return bimg.JPEG
	case "png":
		return bimg.PNG
	case "webp":
		return bimg.WEBP
	case "tiff":
		return bimg.TIFF
	case "gif":
		return bimg.GIF
	case "svg":
		return bimg.SVG
	case "pdf":
		return bimg.PDF
	default:
		return bimg.UNKNOWN
	}
}
