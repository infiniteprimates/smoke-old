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
	go metricsToExpvar(time.Duration(cfg.GetInt(config.MetricsPublishingInterval)) * time.Second)

	go func() {
		ipAndPort := fmt.Sprintf("%s:%d", cfg.GetString(config.MetricsIp), cfg.GetInt(config.MetricsPort))
		log.Info("Starting metrics server on ", ipAndPort)
		log.Fatal(http.ListenAndServe(ipAndPort, nil))
	}()
}

func metricsToExpvar(interval time.Duration) {
	smokeMetrics := expvar.NewMap("smokeMetrics")

	for _ = range time.Tick(interval) {
		metrics.DefaultRegistry.Each(func(name string, i interface{}) {
			switch metric := i.(type) {
			case metrics.Counter:
				smokeMetrics.Add(name, metric.Count())
			case metrics.Gauge:
				smokeMetrics.Add(name, metric.Value())
			case metrics.GaugeFloat64:
				smokeMetrics.AddFloat(name, metric.Value())
			case metrics.Histogram:
				h := metric.Snapshot()
				ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
				smokeMetrics.Add(name+".count", h.Count())
				smokeMetrics.AddFloat(name+".min", float64(metric.Min()))
				smokeMetrics.AddFloat(name+".min", float64(h.Min()))
				smokeMetrics.AddFloat(name+".max", float64(h.Max()))
				smokeMetrics.AddFloat(name+".mean", float64(h.Mean()))
				smokeMetrics.AddFloat(name+".std-dev", float64(h.StdDev()))
				smokeMetrics.AddFloat(name+".50-percentile", float64(ps[0]))
				smokeMetrics.AddFloat(name+".75-percentile", float64(ps[1]))
				smokeMetrics.AddFloat(name+".95-percentile", float64(ps[2]))
				smokeMetrics.AddFloat(name+".99-percentile", float64(ps[3]))
				smokeMetrics.AddFloat(name+".999-percentile", float64(ps[4]))
			case metrics.Meter:
				m := metric.Snapshot()
				smokeMetrics.Add(name+".count", m.Count())
				smokeMetrics.AddFloat(name+".one-minute", float64(m.Rate1()))
				smokeMetrics.AddFloat(name+".five-minute", float64(m.Rate5()))
				smokeMetrics.AddFloat(name+".fifteen-minute", float64((m.Rate15())))
				smokeMetrics.AddFloat(name+".mean", float64(m.RateMean()))
			case metrics.Timer:
				t := metric.Snapshot()
				ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
				smokeMetrics.Add(name+".count", t.Count())
				smokeMetrics.AddFloat(name+".min", float64(t.Min()))
				smokeMetrics.AddFloat(name+".max", float64(t.Max()))
				smokeMetrics.AddFloat(name+".mean", float64(t.Mean()))
				smokeMetrics.AddFloat(name+".std-dev", float64(t.StdDev()))
				smokeMetrics.AddFloat(name+".50-percentile", float64(ps[0]))
				smokeMetrics.AddFloat(name+".75-percentile", float64(ps[1]))
				smokeMetrics.AddFloat(name+".95-percentile", float64(ps[2]))
				smokeMetrics.AddFloat(name+".99-percentile", float64(ps[3]))
				smokeMetrics.AddFloat(name+".999-percentile", float64(ps[4]))
				smokeMetrics.AddFloat(name+".one-minute", float64(t.Rate1()))
				smokeMetrics.AddFloat(name+".five-minute", float64(t.Rate5()))
				smokeMetrics.AddFloat(name+".fifteen-minute", float64((t.Rate15())))
				smokeMetrics.AddFloat(name+".mean-rate", float64(t.RateMean()))
			default:
				panic(fmt.Sprintf("unsupported type for '%s': %T", name, i))
			}
		})
	}
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
