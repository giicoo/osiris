package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/osiris/template-service/internal/config"
	"github.com/osiris/template-service/internal/delivery/restapi"
	"github.com/osiris/template-service/internal/repository"
	"github.com/osiris/template-service/internal/server"
	"github.com/osiris/template-service/internal/services"
	"github.com/osiris/template-service/pkg/logging"
	"github.com/sirupsen/logrus"
)

func RunApp() {
	logging.SetupLogging("template")
	cfg := config.SetupConfig("template")

	repository := repository.NewRepoTemp(cfg)
	services := services.NewServices(cfg, repository)
	controller := restapi.NewController(cfg, services)

	r := restapi.SetupRouter(controller)

	srv := server.NewServer(r.Handler())

	go func() {

		err := srv.StartServer()
		if err != nil {
			switch err {
			case http.ErrServerClosed:

			default:
				logrus.Errorf("start server: %s", err)
				return
			}
		}
	}()
	logrus.Info("Server Start")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	// ShutDown Server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.ShutdownServer(ctx); err != nil {
		logrus.Errorf("shutdown server: %s", err)
		return
	} else {
		logrus.Info("Server stop")
	}
}
