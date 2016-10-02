package server

import (
	"expvar"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/rcrowley/go-metrics"
)

func TestMetrics_metricsToExpvar(t *testing.T) {
	metrics.
	smokeMetrics := expvar.NewMap("smokeMetrics")

	metricsToExpvar(smokeMetrics)

	assert.Equal(t, smokeMetrics.Get("counter", counterValue, "Unexpected counter value"))
}
