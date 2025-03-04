package restapi

import "github.com/gin-gonic/gin"

func SetupRouter(controller *Controller) *gin.Engine {
	r := gin.Default()

	r.GET("/", controller.Auth)
	r.GET("/callback", controller.Callback)

	return r
}
