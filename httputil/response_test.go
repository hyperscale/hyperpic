// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package httputil

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()

	JSON(w, http.StatusCreated, map[string]string{
		"foo": "bar",
	})

	assert.Equal(t, 201, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.HeaderMap.Get("Content-Type"))
	assert.Equal(t, "nosniff", w.HeaderMap.Get("X-Content-Type-Options"))
	assert.JSONEq(t, `{"foo":"bar"}`, w.Body.String())
}

func TestJSONFail(t *testing.T) {
	w := &ResponseWriterFailMock{}

	JSON(w, http.StatusCreated, map[string]string{
		"foo": "bar",
	})

	assert.Equal(t, 201, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.HeaderMap.Get("Content-Type"))
	assert.Equal(t, "nosniff", w.HeaderMap.Get("X-Content-Type-Options"))
}

func TestFailure(t *testing.T) {
	w := httptest.NewRecorder()

	Failure(w, http.StatusInternalServerError, ErrorMessage{
		Code:    1337,
		Message: "user_message",
	})

	assert.Equal(t, 500, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.HeaderMap.Get("Content-Type"))
	assert.Equal(t, "nosniff", w.HeaderMap.Get("X-Content-Type-Options"))
	assert.JSONEq(t, `{"error":{"code":1337,"message":"user_message"}}`, w.Body.String())
}

func TestFailureFromError(t *testing.T) {
	w := httptest.NewRecorder()

	FailureFromError(w, http.StatusInternalServerError, errors.New("dev_message"))

	assert.Equal(t, 500, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.HeaderMap.Get("Content-Type"))
	assert.Equal(t, "nosniff", w.HeaderMap.Get("X-Content-Type-Options"))
	assert.JSONEq(t, `{"error":{"code":500,"message":"dev_message"}}`, w.Body.String())
}

func TestNotFoundFailure(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/users", nil)
	w := httptest.NewRecorder()

	NotFoundFailure(w, req)

	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.HeaderMap.Get("Content-Type"))
	assert.Equal(t, "nosniff", w.HeaderMap.Get("X-Content-Type-Options"))
	assert.JSONEq(t, `{"error":{"code":404,"message":"No route found for \\\"GET /users\\\""}}`, w.Body.String())
}
