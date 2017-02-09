// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

func main() {
	config := NewConfiguration()

	// for dev
	config.Image.Path = "." + config.Image.Path
	config.Server.Port = 8574

	server := NewServer(config)

	server.ListenAndServe()
}
