package middleware

import (
	"github.com/gin-gonic/gin"
)

func InitMiddleware(service []interface{}) gin.HandlerFunc {
	return func(context *gin.Context) {
		// 将实例存在gin.Keys中
		context.Keys = make(map[string]interface{})
		context.Keys["user"] = service[0]
		context.Keys["task"] = service[1]
		context.Next()
	}
}
