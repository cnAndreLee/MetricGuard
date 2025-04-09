package routes

import (
	"github.com/cnAndreLee/MetricGuard/controller"

	"github.com/gin-gonic/gin"
)

func GetSmnRoutes(route *gin.RouterGroup) {

	smn := route.Group("/smn")

	smn.POST("/alert", controller.Alert)
}
