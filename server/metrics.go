package server

import (
	"expvar"
	"fmt"
	"net/http"
	_ "net/http/pprof"
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
			setFloat(expvarMap, name+".min", float64(snapshot.Min()))
			setFloat(expvarMap, name+".max", float64(snapshot.Max()))
			setFloat(expvarMap, name+".mean", snapshot.Mean())
			setFloat(expvarMap, name+".std-dev", snapshot.StdDev())
			setFloat(expvarMap, name+".p50", percentiles[0])
			setFloat(expvarMap, name+".p75", percentiles[1])
			setFloat(expvarMap, name+".p95", percentiles[2])
			setFloat(expvarMap, name+".p99", percentiles[3])
			setFloat(expvarMap, name+".p999", percentiles[4])
		case metrics.Meter:
			m := metric.Snapshot()
			expvarMap.Add(name+".count", m.Count())
			expvarMap.AddFloat(name+".one-minute", float64(m.Rate1()))
			expvarMap.AddFloat(name+".five-minute", float64(m.Rate5()))
			expvarMap.AddFloat(name+".fifteen-minute", float64((m.Rate15())))
			expvarMap.AddFloat(name+".mean", float64(m.RateMean()))
		case metrics.Timer:
			t := metric.Snapshot()
			ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			expvarMap.Add(name+".count", t.Count())
			expvarMap.AddFloat(name+".min", float64(t.Min()))
			expvarMap.AddFloat(name+".max", float64(t.Max()))
			expvarMap.AddFloat(name+".mean", float64(t.Mean()))
			expvarMap.AddFloat(name+".std-dev", float64(t.StdDev()))
			expvarMap.AddFloat(name+".50-percentile", float64(ps[0]))
			expvarMap.AddFloat(name+".75-percentile", float64(ps[1]))
			expvarMap.AddFloat(name+".95-percentile", float64(ps[2]))
			expvarMap.AddFloat(name+".99-percentile", float64(ps[3]))
			expvarMap.AddFloat(name+".999-percentile", float64(ps[4]))
			expvarMap.AddFloat(name+".one-minute", float64(t.Rate1()))
			expvarMap.AddFloat(name+".five-minute", float64(t.Rate5()))
			expvarMap.AddFloat(name+".fifteen-minute", float64((t.Rate15())))
			expvarMap.AddFloat(name+".mean-rate", float64(t.RateMean()))
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
