// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	prometheus.MustRegister(CacheHit)
	prometheus.MustRegister(CacheMiss)
	prometheus.MustRegister(ImageDeliveredBytes)
	prometheus.MustRegister(ImageReceivedBytes)
}

// CacheHit counter
var CacheHit = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "cache_hit_total",
		Help: "The count of cache hit.",
	},
	[]string{},
)

// CacheMiss counter
var CacheMiss = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "cache_miss_total",
		Help: "The count of cache miss.",
	},
	[]string{},
)

// ImageDeliveredBytes counter
var ImageDeliveredBytes = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "image_delivered_bytes",
		Help: "The bytes of image delivered.",
	},
	[]string{},
)

// ImageReceivedBytes counter
var ImageReceivedBytes = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "image_received_bytes",
		Help: "The bytes of image received.",
	},
	[]string{},
)
