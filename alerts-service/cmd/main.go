package main

import (
	"github.com/giicoo/osiris/alerts-service/docs"
	"github.com/giicoo/osiris/alerts-service/internal/app"
)

//	@title		Osiris Alerts Service API
//	@version	1.0

// @host		giicoo.ru
// @BasePath /api/alerts-service
func main() {
	docs.SwaggerInfo.BasePath = "/api/alerts-service"
	app.RunApp()
}
