package middleware

import (
	dgerr "github.com/darwinOrg/go-common/enums/error"
	"github.com/gin-gonic/gin"
	"net/http"
)

var tooManyRequestError = &dgerr.DgError{
	Code:    http.StatusBadRequest,
	Message: "请求太频繁",
}

func SkippedPathPrefixes(c *gin.Context, prefixes ...string) bool {
	if len(prefixes) == 0 {
		return false
	}

	path := c.Request.URL.Path
	pathLen := len(path)
	for _, p := range prefixes {
		if pl := len(p); pathLen >= pl && path[:pl] == p {
			return true
		}
	}
	return false
}

func AllowedPathPrefixes(c *gin.Context, prefixes ...string) bool {
	if len(prefixes) == 0 {
		return true
	}

	path := c.Request.URL.Path
	pathLen := len(path)
	for _, p := range prefixes {
		if pl := len(p); pathLen >= pl && path[:pl] == p {
			return true
		}
	}
	return false
}

func Empty() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
