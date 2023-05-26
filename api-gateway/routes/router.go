package routes

import (
	"main/ToDoList_grpc/api-gateway/inner/handler"
	"main/ToDoList_grpc/api-gateway/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(service ...interface{}) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors(), middleware.InitMiddleware(service))
	v1 := r.Group("/api/v1")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, "SUCCESS") //ping
		})
		v1.POST("/user/register", handler.UserRegister)
		v1.POST("/user/login", handler.UserLogin)

	}

	// 需要登录保护
	authed := v1.Group("/")
	authed.Use(middleware.JWT())
	{
		// 任务模块
		authed.GET("/task/show", handler.ListTask)
		authed.POST("/task/create", handler.CreateTask)
		authed.PUT("/task/update", handler.UpdateTask)
		authed.DELETE("/task/delete", handler.DeleteTask)
	}

	return r
}
