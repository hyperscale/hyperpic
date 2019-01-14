// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package image

import (
	"strings"

	"github.com/h2non/bimg"
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

	return IsFormatSupported(format)
}

// IsFormatSupported returns true if the image format is supported by bimg.
func IsFormatSupported(format string) bool {
	// Some payloads may expose the MIME type for SVG as text/xml
	if format == "xml" {
		format = "svg"
	}

	if format == "jpg" {
		format = "jpeg"
	}

	return bimg.IsTypeNameSupported(format)
}

// ExtensionToType returns the image type based on the given image type alias.
func ExtensionToType(name string) bimg.ImageType {
	ext := strings.ToLower(name)

	switch ext {
	case "jpeg", "jpg":
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
