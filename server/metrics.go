package server

import (
	"expvar"
	"fmt"
	"net/http"
	_ "net/http/pprof" // Needed to enable pprof
	"time"

	"github.com/infiniteprimates/smoke/config"
	"github.com/labstack/echo"
	"github.com/labstack/echo/log"
	"github.com/rcrowley/go-metrics"
)

func startMetricsServer(cfg config.Config, log log.Logger) {
	// start publishing metrics into expvars
	go gatherMetrics(time.Duration(cfg.GetInt(config.MetricsPublishingInterval)) * time.Second)

	go func() {
		ipAndPort := fmt.Sprintf("%s:%d", cfg.GetString(config.MetricsIp), cfg.GetInt(config.MetricsPort))
		log.Info("Starting metrics server on ", ipAndPort)
		log.Fatal(http.ListenAndServe(ipAndPort, nil))
	}()
}

func gatherMetrics(interval time.Duration) {
	smokeMetrics := expvar.NewMap("smokeMetrics")

	for _ = range time.Tick(interval) {
		metricsToExpvar(smokeMetrics)
	}
}

func metricsToExpvar(expvarMap *expvar.Map) {
	metrics.DefaultRegistry.Each(func(name string, i interface{}) {
		switch metric := i.(type) {
		case metrics.Counter:
			setInt(expvarMap, name, metric.Count())
		case metrics.Gauge:
			setInt(expvarMap, name, metric.Value())
		case metrics.GaugeFloat64:
			var gauge expvar.Float
			gauge.Set(metric.Value())
			expvarMap.Set(name, &gauge)
		case metrics.Histogram:
			snapshot := metric.Snapshot()
			percentiles := snapshot.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			setInt(expvarMap, name+".count", snapshot.Count())
			setInt(expvarMap, name+".min", snapshot.Min())
			setInt(expvarMap, name+".max", snapshot.Max())
			setFloat(expvarMap, name+".mean", snapshot.Mean())
			setFloat(expvarMap, name+".std-dev", snapshot.StdDev())
			setFloat(expvarMap, name+".p50", percentiles[0])
			setFloat(expvarMap, name+".p75", percentiles[1])
			setFloat(expvarMap, name+".p95", percentiles[2])
			setFloat(expvarMap, name+".p99", percentiles[3])
			setFloat(expvarMap, name+".p999", percentiles[4])
		case metrics.Meter:
			snapshot := metric.Snapshot()
			setInt(expvarMap, name+".count", snapshot.Count())
			setFloat(expvarMap, name+".1min", snapshot.Rate1())
			setFloat(expvarMap, name+".5min", snapshot.Rate5())
			setFloat(expvarMap, name+".15min", snapshot.Rate15())
			setFloat(expvarMap, name+".mean-rate", snapshot.RateMean())
		case metrics.Timer:
			snapshot := metric.Snapshot()
			percentiles := snapshot.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			setInt(expvarMap, name+".count", snapshot.Count())
			setInt(expvarMap, name+".min", snapshot.Min())
			setInt(expvarMap, name+".max", snapshot.Max())
			setFloat(expvarMap, name+".mean", snapshot.Mean())
			setFloat(expvarMap, name+".std-dev", snapshot.StdDev())
			setFloat(expvarMap, name+".p50", percentiles[0])
			setFloat(expvarMap, name+".p75", percentiles[1])
			setFloat(expvarMap, name+".p95", percentiles[2])
			setFloat(expvarMap, name+".p99", percentiles[3])
			setFloat(expvarMap, name+".p999", percentiles[4])
			setFloat(expvarMap, name+".1min", snapshot.Rate1())
			setFloat(expvarMap, name+".5min", snapshot.Rate5())
			setFloat(expvarMap, name+".15min", snapshot.Rate15())
			setFloat(expvarMap, name+".mean-rate", snapshot.RateMean())
		default:
			panic(fmt.Sprintf("unsupported type for '%s': %T", name, i))
		}
	})
}

func setInt(expvarMap *expvar.Map, name string, value int64) {
	var expvarInt expvar.Int
	expvarInt.Set(value)
	expvarMap.Set(name, &expvarInt)
}

func setFloat(expvarMap *expvar.Map, name string, value float64) {
	var expvarFloat expvar.Float
	expvarFloat.Set(value)
	expvarMap.Set(name, &expvarFloat)
}

func metricsMiddleware(name string) echo.MiddlewareFunc {
	timer := metrics.NewTimer()

	metrics.Register(fmt.Sprintf("api.%s.timer", name), timer)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			defer timer.UpdateSince(start)
			return next(c)
		}
	}
}
