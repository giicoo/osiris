package restapi

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(controller *Controller) *gin.Engine {
	r := gin.Default()
	r.Use(AuthUser())

	r.POST("/create/alert", controller.CreateAlert)
	r.POST("/create/type", controller.CreateType)
	r.POST("/stop/alert", controller.StopAlert)

	r.GET("/get/alert/:id", controller.GetAlert)
	r.GET("/get/type/:id", controller.GetType)

	r.GET("/get/alerts", controller.GetAlerts)
	r.GET("/get/types", controller.GetTypes)

	r.DELETE("/delete/type", controller.DeleteType)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}
