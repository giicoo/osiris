package restapi

import "github.com/gin-gonic/gin"

func SetupRouter(controller *Controller) *gin.Engine {
	r := gin.Default()
	// TODO: auth with websocket
	r.GET("/ws", controller.WsConnect)
	return r
}
