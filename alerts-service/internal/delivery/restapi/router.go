package restapi

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(controller *Controller) *gin.Engine {
	r := gin.Default()

	{
		gr := r.Group("/")
		gr.Use(AuthUser())

		gr.POST("/create/alert", controller.CreateAlert)
		gr.POST("/create/type", controller.CreateType)
		gr.POST("/stop/alert", controller.StopAlert)

		gr.GET("/get/alert/:id", controller.GetAlert)
		gr.GET("/get/type/:id", controller.GetType)

		gr.GET("/get/alerts", controller.GetAlerts)
		gr.GET("/get/types", controller.GetTypes)

		gr.DELETE("/delete/type", controller.DeleteType)
	}
	r.Static("/static", "./dist")

	return r
}
