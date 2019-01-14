// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package httputil

// ErrorMessageInterface interface
type ErrorMessageInterface interface {
	GetCode() int
	GetMessage() string
	Error() string
}

// ErrorMessage struct
type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error message
func (e ErrorMessage) Error() string {
	return e.Message
}

// GetCode error
func (e ErrorMessage) GetCode() int {
	return e.Code
}

// GetMessage error
func (e ErrorMessage) GetMessage() string {
	return e.Message
}

// ErrorResponse struct
type ErrorResponse struct {
	Error ErrorMessage `json:"error"`
}
