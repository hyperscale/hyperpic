// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"github.com/h2non/bimg"
)

func TestExtractImageTypeFromMime(t *testing.T) {
	files := []struct {
		mime     string
		expected string
	}{
		{"image/jpeg", "jpeg"},
		{"/png", "png"},
		{"png", ""},
		{"multipart/form-data; encoding=utf-8", "form-data"},
		{"", ""},
	}

	for _, file := range files {
		if ExtractImageTypeFromMime(file.mime) != file.expected {
			t.Fatalf("Invalid mime type: %s != %s", file.mime, file.expected)
		}
	}
}

func TestIsImageTypeSupported(t *testing.T) {
	files := []struct {
		name     string
		expected bool
	}{
		{"image/jpeg", true},
		{"image/png", true},
		{"image/webp", bimg.IsTypeSupported(bimg.WEBP)},
		{"IMAGE/JPEG", true},
		{"png", false},
		{"multipart/form-data; encoding=utf-8", false},
		{"application/json", false},
		{"image/gif", bimg.IsTypeSupported(bimg.GIF)},
		{"image/svg+xml", bimg.IsTypeSupported(bimg.SVG)},
		{"image/svg", bimg.IsTypeSupported(bimg.SVG)},
		{"image/tiff", bimg.IsTypeSupported(bimg.TIFF)},
		{"application/pdf", bimg.IsTypeSupported(bimg.PDF)},
		{"text/plain", false},
		{"blablabla", false},
		{"", false},
	}

	for _, file := range files {
		if IsImageMimeTypeSupported(file.name) != file.expected {
			t.Fatalf("Invalid type: %s != %t", file.name, file.expected)
		}
	}
}

func TestIsFormatSupported(t *testing.T) {
	files := []struct {
		name     string
		expected bool
	}{
		{"jpeg", true},
		{"jpg", true},
		{"png", true},
		{"webp", bimg.IsTypeSupported(bimg.WEBP)},
		{"JPEG", false},
		{"json", false},
		{"gif", bimg.IsTypeSupported(bimg.GIF)},
		{"xml", bimg.IsTypeSupported(bimg.SVG)},
		{"svg", bimg.IsTypeSupported(bimg.SVG)},
		{"tiff", bimg.IsTypeSupported(bimg.TIFF)},
	}

	for _, file := range files {
		if IsFormatSupported(file.name) != file.expected {
			t.Fatalf("Invalid type: %s != %t", file.name, file.expected)
		}
	}
}

func TestImageType(t *testing.T) {
	files := []struct {
		name     string
		expected bimg.ImageType
	}{
		{"jpeg", bimg.JPEG},
		{"jpg", bimg.JPEG},
		{"png", bimg.PNG},
		{"webp", bimg.WEBP},
		{"tiff", bimg.TIFF},
		{"gif", bimg.GIF},
		{"svg", bimg.SVG},
		{"pdf", bimg.PDF},
		{"multipart/form-data; encoding=utf-8", bimg.UNKNOWN},
		{"json", bimg.UNKNOWN},
		{"text", bimg.UNKNOWN},
		{"blablabla", bimg.UNKNOWN},
		{"", bimg.UNKNOWN},
	}

	for _, file := range files {
		if ImageType(file.name) != file.expected {
			t.Fatalf("Invalid type: %s != %v", file.name, file.expected)
		}
	}
}

func TestGetImageMimeType(t *testing.T) {
	files := []struct {
		name     bimg.ImageType
		expected string
	}{
		{bimg.JPEG, "image/jpeg"},
		{bimg.PNG, "image/png"},
		{bimg.WEBP, "image/webp"},
		{bimg.TIFF, "image/tiff"},
		{bimg.GIF, "image/gif"},
		{bimg.PDF, "application/pdf"},
		{bimg.SVG, "image/svg+xml"},
		{bimg.UNKNOWN, "image/jpeg"},
	}

	for _, file := range files {
		if GetImageMimeType(file.name) != file.expected {
			t.Fatalf("Invalid type: %v != %v", file.name, file.expected)
		}
	}
}
