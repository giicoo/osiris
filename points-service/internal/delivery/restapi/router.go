package restapi

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(controller *Controller) *gin.Engine {
	r := gin.Default()

	r.POST("/create/point", controller.CreatePoint)
	r.GET("/get/point/:id", controller.GetPoint)
	r.GET("/get/points/:user_id", controller.GetPoints)
	r.PUT("/update/point/location", controller.UpdateLocation)
	r.PUT("/update/point/title", controller.UpdateTitle)
	r.PUT("/update/point/radius", controller.UpdateRadius)
	r.DELETE("/delete/point", controller.DeletePoint)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}
