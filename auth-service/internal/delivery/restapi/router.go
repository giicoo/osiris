package restapi

import "github.com/gin-gonic/gin"

func SetupRouter(controller *Controller) *gin.Engine {
	r := gin.Default()

	r.GET("/auth", controller.Auth)
	r.GET("/callback", controller.Callback)
	r.GET("/check-user", controller.CheckUser)

	return r
}
