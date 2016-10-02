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
	gauge.Update(220)

	var smokeMetrics expvar.Map
	smokeMetrics.Init()

	metricsToExpvar(&smokeMetrics)
	// a second time to make sure things aren't incorrectly accumulated
	metricsToExpvar(&smokeMetrics)

	assert.Equal(t, "220", smokeMetrics.Get("gauge").String(), "Unexpected gauge value.")
}

func TestMetrics_metricsToExpvar_GaugeFloat64(t *testing.T) {
	gaugefloat64 := metrics.NewGaugeFloat64()
	metrics.Register("gaugefloat64", gaugefloat64)
	gaugefloat64.Update(.223)

	var smokeMetrics expvar.Map
	smokeMetrics.Init()

	metricsToExpvar(&smokeMetrics)
	// a second time to make sure things aren't incorrectly accumulated
	metricsToExpvar(&smokeMetrics)

	assert.Equal(t, "0.223", smokeMetrics.Get("gaugefloat64").String(), "Unexpected gaugefloat64 value.")
}

func TestMetrics_metricsToExpvar_Histogram(t *testing.T) {
	sample := metrics.NewUniformSample(4)
	histogram := metrics.NewHistogram(sample)
	metrics.Register("histogram", histogram)
	histogram.Update(223)
	histogram.Update(260)
	histogram.Update(270)
	histogram.Update(308)
	histogram.Update(357)

	var smokeMetrics expvar.Map
	smokeMetrics.Init()

	metricsToExpvar(&smokeMetrics)
	// a second time to make sure things aren't incorrectly accumulated
	metricsToExpvar(&smokeMetrics)

	assert.Equal(t, "5", smokeMetrics.Get("histogram.count").String(), "Unexpected histogram count value.")
	assert.Equal(t, "260", smokeMetrics.Get("histogram.min").String(), "Unexpected histogram min value.")
	assert.Equal(t, "357", smokeMetrics.Get("histogram.max").String(), "Unexpected histogram max value.")
	assert.Equal(t, "298.75", smokeMetrics.Get("histogram.mean").String(), "Unexpected histogram mean value.")
	assert.Equal(t, "38.10101704679286", smokeMetrics.Get("histogram.std-dev").String(), "Unexpected histogram std-dev value.")
	assert.Equal(t, "289", smokeMetrics.Get("histogram.p50").String(), "Unexpected histogram p50 value.")
	assert.Equal(t, "344.75", smokeMetrics.Get("histogram.p75").String(), "Unexpected histogram p75 value.")
	assert.Equal(t, "357", smokeMetrics.Get("histogram.p95").String(), "Unexpected histogram p95 value.")
	assert.Equal(t, "357", smokeMetrics.Get("histogram.p99").String(), "Unexpected histogram p99 value.")
	assert.Equal(t, "357", smokeMetrics.Get("histogram.p999").String(), "Unexpected histogram p999 value.")
}
