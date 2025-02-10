package restapi

import "github.com/gin-gonic/gin"

func SetupRouter(controller *Controller) *gin.Engine {
	r := gin.Default()

	return r
}
