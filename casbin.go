package dgwe

import (
	"github.com/casbin/casbin/v2"
	dgerr "github.com/darwinOrg/go-common/enums/error"
	"github.com/darwinOrg/go-common/result"
	"github.com/darwinOrg/go-web/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CasbinConfig struct {
	AllowedPathPrefixes []string
	SkippedPathPrefixes []string
	Skipper             func(c *gin.Context) bool
	GetEnforcer         func(c *gin.Context) *casbin.Enforcer
	GetSubjects         func(c *gin.Context) []string
}

func CasbinWithConfig(config CasbinConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !middleware.AllowedPathPrefixes(c, config.AllowedPathPrefixes...) ||
			middleware.SkippedPathPrefixes(c, config.SkippedPathPrefixes...) ||
			(config.Skipper != nil && config.Skipper(c)) {
			c.Next()
			return
		}

		enforcer := config.GetEnforcer(c)
		if enforcer == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, result.SimpleFailByError(dgerr.SYSTEM_ERROR))
			return
		}

		for _, sub := range config.GetSubjects(c) {
			if b, err := enforcer.Enforce(sub, c.Request.URL.Path, c.Request.Method); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, result.SimpleFailByError(err))
				return
			} else if b {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusOK, result.SimpleFailByError(dgerr.NO_PERMISSION))
	}
}
