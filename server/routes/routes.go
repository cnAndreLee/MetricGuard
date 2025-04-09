package routes

import (
	"github.com/cnAndreLee/MetricGuard/middleware"

	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {

	// 解决跨域问题
	r.Use(middleware.CORSMiddleware())

	v1 := r.Group("/api/v1")

	{
		GetSmnRoutes(v1)
	}

	v1.GET("ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return r

}
