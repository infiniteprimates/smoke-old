package server

import (
	"expvar"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/rcrowley/go-metrics"
)

func TestMetrics_metricsToExpvar_Counter(t *testing.T) {
	counter := metrics.NewCounter()
	metrics.Register("counter", counter)
	counter.Inc(22)

	var smokeMetrics expvar.Map
	smokeMetrics.Init()

	metricsToExpvar(&smokeMetrics)
	// a second time to make sure things aren't incorrectly accumulated
	metricsToExpvar(&smokeMetrics)

	assert.Equal(t, "22", smokeMetrics.Get("counter").String(), "Unexpected counter value.")
}

func TestMetrics_metricsToExpvar_Gauge(t *testing.T) {
	gauge := metrics.NewGauge()
	metrics.Register("gauge", gauge)
	gauge.Update(223)

	var smokeMetrics expvar.Map
	smokeMetrics.Init()

	metricsToExpvar(&smokeMetrics)
	// a second time to make sure things aren't incorrectly accumulated
	metricsToExpvar(&smokeMetrics)

	assert.Equal(t, "223", smokeMetrics.Get("gauge").String(), "Unexpected gauge value.")
}

func TestMetrics_metricsToExpvar_GaugeFloat64(t *testing.T) {
	gaugefloat64 := metrics.NewGaugeFloat64()
	metrics.Register("gaugefloat64", gaugefloat64)
	gaugefloat64.Update(3.08)

	var smokeMetrics expvar.Map
	smokeMetrics.Init()

	metricsToExpvar(&smokeMetrics)
	// a second time to make sure things aren't incorrectly accumulated
	metricsToExpvar(&smokeMetrics)

	assert.Equal(t, "3.08", smokeMetrics.Get("gaugefloat64").String(), "Unexpected gaugefloat64 value.")
}
