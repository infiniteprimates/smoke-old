package server

import (
	"expvar"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/test"
	"github.com/rcrowley/go-metrics"
	"github.com/stretchr/testify/assert"
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

func TestMetrics_metricsToExpvar_Meter(t *testing.T) {
	meter := metrics.NewMeter()
	metrics.Register("meter", meter)
	meter.Mark(410)
	meter.Mark(44)
	meter.Mark(45)
	meter.Mark(454)

	var smokeMetrics expvar.Map
	smokeMetrics.Init()

	metricsToExpvar(&smokeMetrics)
	// a second time to make sure things aren't incorrectly accumulated
	metricsToExpvar(&smokeMetrics)

	assert.Equal(t, "953", smokeMetrics.Get("meter.count").String(), "Unexpected meter count value.")
	assert.Equal(t, "0", smokeMetrics.Get("meter.1min").String(), "Unexpected meter 1min value.")
	assert.Equal(t, "0", smokeMetrics.Get("meter.5min").String(), "Unexpected meter 5min value.")
	assert.Equal(t, "0", smokeMetrics.Get("meter.15min").String(), "Unexpected meter 15min value.")
	assert.NotZero(t, smokeMetrics.Get("meter.mean-rate").String(), "Unexpected meter 15min value.")
}

func TestMetrics_metricsToExpvar_Timer(t *testing.T) {
	timer := metrics.NewTimer()
	metrics.Register("timer", timer)

	timer.Update(2 * time.Microsecond)
	timer.Update(3 * time.Microsecond)
	timer.Update(5 * time.Microsecond)
	timer.Update(7 * time.Microsecond)
	timer.Update(11 * time.Microsecond)

	var smokeMetrics expvar.Map
	smokeMetrics.Init()

	metricsToExpvar(&smokeMetrics)
	// a second time to make sure things aren't incorrectly accumulated
	metricsToExpvar(&smokeMetrics)

	assert.Equal(t, "5", smokeMetrics.Get("timer.count").String(), "Unexpected timer count value.")
	assert.Equal(t, "2000", smokeMetrics.Get("timer.min").String(), "Unexpected timer min value.")
	assert.Equal(t, "11000", smokeMetrics.Get("timer.max").String(), "Unexpected timer max value.")
	assert.Equal(t, "5600", smokeMetrics.Get("timer.mean").String(), "Unexpected timer mean value.")
	assert.Equal(t, "3200", smokeMetrics.Get("timer.std-dev").String(), "Unexpected timer std-dev value.")
	assert.Equal(t, "5000", smokeMetrics.Get("timer.p50").String(), "Unexpected timer p50 value.")
	assert.Equal(t, "9000", smokeMetrics.Get("timer.p75").String(), "Unexpected timer p75 value.")
	assert.Equal(t, "11000", smokeMetrics.Get("timer.p95").String(), "Unexpected timer p95 value.")
	assert.Equal(t, "11000", smokeMetrics.Get("timer.p99").String(), "Unexpected timer p99 value.")
	assert.Equal(t, "11000", smokeMetrics.Get("timer.p999").String(), "Unexpected timer p999 value.")
	assert.Equal(t, "0", smokeMetrics.Get("timer.1min").String(), "Unexpected timer 1min value.")
	assert.Equal(t, "0", smokeMetrics.Get("timer.5min").String(), "Unexpected timer 5min value.")
	assert.Equal(t, "0", smokeMetrics.Get("timer.15min").String(), "Unexpected timer 15min value.")
	assert.NotZero(t, "578770.6910522052", smokeMetrics.Get("timer.mean-rate").String(), "Unexpected timer mean-rate value.")
}

func TestMetrics_metricsMiddleware(t *testing.T) {
	e := echo.New()
	req := test.NewRequest(echo.GET, "/", strings.NewReader(""))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)

	mw := metricsMiddleware("rest")
	err := mw(func(ctx echo.Context) error { return c.String(http.StatusOK, "OK") })(c)

	if assert.NoError(t, err, "An error occurred invoking handler.") {
		timer := metrics.DefaultRegistry.Get("api.rest.timer").(metrics.Timer)
		assert.NotZero(t, timer.Min(), "Timer duration was zero.")

	}
}
