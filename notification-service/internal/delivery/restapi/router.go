package restapi

import "github.com/gin-gonic/gin"

func SetupRouter(controller *Controller) *gin.Engine {
	r := gin.Default()
	r.GET("/ws/:id", controller.WsConnect)
	return r
}
