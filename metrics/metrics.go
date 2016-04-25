package metrics

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rcrowley/go-metrics"
)

func MetricsHandler(name string) gin.HandlerFunc {
	counter := metrics.NewCounter()
	timer := metrics.NewTimer()

	metrics.Register(fmt.Sprintf("api.%s.counter", name), counter)
	metrics.Register(fmt.Sprintf("api.%s.timer", name), timer)

	return func(ctx *gin.Context) {
		counter.Inc(1)
		start := time.Now()
		ctx.Next()
		timer.UpdateSince(start)
	}
}
