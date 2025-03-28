package restapi

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(controller *Controller) *gin.Engine {
	r := gin.Default()
	r.POST("/registration", controller.CreateUser)
	r.GET("/get/user/:id", controller.GetUser)

	r.POST("/login", controller.Login)
	r.GET("/auth/:session_id", controller.Auth)

	return r
}
