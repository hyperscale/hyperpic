// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestMetrics(t *testing.T) {
	assert.True(t, prometheus.Unregister(CacheHit))
	assert.True(t, prometheus.Unregister(CacheMiss))
	assert.True(t, prometheus.Unregister(ImageDeliveredBytes))
	assert.True(t, prometheus.Unregister(ImageReceivedBytes))
}
