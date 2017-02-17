package httputil

import (
	"net/http"
	"reflect"
)

// ResponseStatus returns the HTTP response status.
// Remember that the status is only set by the server after WriteHeader has been called.
func ResponseStatus(w http.ResponseWriter) int {
	return int(httpResponseStruct(reflect.ValueOf(w)).FieldByName("status").Int())
}

// httpResponseStruct returns the response structure after going trough all the intermediary response writers.
func httpResponseStruct(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Type().String() == "http.response" {
		return v
	}

	return httpResponseStruct(v.FieldByName("ResponseWriter").Elem())
}
