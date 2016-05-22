package server

import (
	"fmt"
	"time"

	"github.com/labstack/echo"
	"github.com/rcrowley/go-metrics"
)

func metricsHandler(name string) echo.MiddlewareFunc {
	counter := metrics.NewCounter()
	timer := metrics.NewTimer()

	metrics.Register(fmt.Sprintf("api.%s.counter", name), counter)
	metrics.Register(fmt.Sprintf("api.%s.timer", name), timer)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			counter.Inc(1)
			start := time.Now()
			defer timer.UpdateSince(start)
			return next(c)
		}
	}
}
