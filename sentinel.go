package dgwe

import (
	"github.com/alibaba/sentinel-golang/core/system"
	sentinel "github.com/alibaba/sentinel-golang/pkg/adapters/gin"
	dgerr "github.com/darwinOrg/go-common/enums/error"
	"github.com/darwinOrg/go-common/result"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// SentinelByQPS 按QPS限流
func SentinelByQPS(triggerCount float64) gin.HandlerFunc {
	if _, err := system.LoadRules([]*system.Rule{
		{
			MetricType:   system.InboundQPS,
			TriggerCount: triggerCount,
			Strategy:     system.BBR,
		},
	}); err != nil {
		log.Fatalf("Unexpected error: %+v", err)
	}

	return sentinel.SentinelMiddleware(
		sentinel.WithBlockFallback(func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusBadRequest, result.SimpleFailByError(dgerr.TOO_MANY_REQUEST))
		}),
	)
}
