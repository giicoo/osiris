package restapi

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(controller *Controller) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Можно указать конкретный домен, например, "http://localhost:3000"
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.POST("/registration", controller.CreateUser)
	r.GET("/get/user/:id", controller.GetUser)

	r.POST("/login", controller.Login)
	r.POST("/logout", controller.Logout)

	r.GET("/auth/:session_id", controller.Auth)

	return r
}
