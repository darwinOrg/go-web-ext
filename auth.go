package dgwe

import (
	"github.com/darwinOrg/go-common/result"
	"github.com/darwinOrg/go-web/middleware"
	"github.com/darwinOrg/go-web/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthConfig struct {
	AllowedPathPrefixes []string
	SkippedPathPrefixes []string
	Skipper             func(c *gin.Context) bool
	ParseUserIdHandler  func(c *gin.Context) (int64, error)
}

func AuthWithConfig(config AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !middleware.AllowedPathPrefixes(c, config.AllowedPathPrefixes...) ||
			middleware.SkippedPathPrefixes(c, config.SkippedPathPrefixes...) ||
			(config.Skipper != nil && config.Skipper(c)) {
			c.Next()
			return
		}

		userId, err := config.ParseUserIdHandler(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, result.SimpleFailByError(err))
			return
		}

		ctx := utils.GetDgContext(c)
		ctx.UserId = userId

		c.Next()
	}
}
