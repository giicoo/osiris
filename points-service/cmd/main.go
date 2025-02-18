package main

import (
	"github.com/giicoo/osiris/points-service/docs"
	"github.com/giicoo/osiris/points-service/internal/app"
)

//	@title		Osiris Points Service API
//	@version	1.0

// @host		giicoo.ru
// @BasePath /api/points-service
func main() {
	docs.SwaggerInfo.BasePath = "/api/points-service"
	app.RunApp()
}
