package main

import (
	"github.com/giicoo/osiris/alerts-service/internal/app"
)

// @Version 1.0.0
// @Title Alert Serive API
// @LicenseName MIT
// @LicenseURL https://en.wikipedia.org/wiki/MIT_License

// @Server https://giicoo.ru/api/alerts-service PROD
// @Server / DEV

// @Security AuthorizationHeader read write
// @SecurityScheme AuthorizationHeader http bearer Input your token
func main() {
	app.RunApp()
}
