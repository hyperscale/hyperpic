// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package httputil

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponseStatus(t *testing.T) {
	statusWant := http.StatusForbidden
	var statusGot, customStatusGot int

	type CustomResponseWriter struct {
		http.ResponseWriter
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, http.StatusText(statusWant), statusWant)
		statusGot = ResponseStatus(w)
		customStatusGot = ResponseStatus(CustomResponseWriter{w})
	}))
	defer ts.Close()

	if _, err := http.Get(ts.URL); err != nil {
		t.Fatal(err)
	}

	if statusWant != statusGot {
		t.Errorf("http.ResponseWriter: want %d, got %d", statusWant, statusGot)
	}

	if customStatusGot != statusGot {
		t.Errorf("CustomResponseWriter: want %d, got %d", customStatusGot, statusGot)
	}
}
