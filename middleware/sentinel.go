package middleware

import (
	"github.com/alibaba/sentinel-golang/core/system"
	sentinel "github.com/alibaba/sentinel-golang/pkg/adapters/gin"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Sentinel 限流
func Sentinel() gin.HandlerFunc {
	if _, err := system.LoadRules([]*system.Rule{
		{
			MetricType:   system.InboundQPS,
			TriggerCount: 200,
			Strategy:     system.BBR,
		},
	}); err != nil {
		log.Fatalf("Unexpected error: %+v", err)
	}

	return sentinel.SentinelMiddleware(
		sentinel.WithBlockFallback(func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusBadRequest, tooManyRequestError)
		}),
	)
}
