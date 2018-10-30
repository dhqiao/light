package middleware

import (
	"github.com/gin-gonic/gin"
	"light/http/global"
)

func SetGlobalData() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.Writer.Header().Get("X-Request-Id")
		if requestId != "" {
			global.SetRequestId(requestId)
		}
		c.Next()
	}
}


