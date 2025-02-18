package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/giicoo/osiris/points-service/internal/config"
	"github.com/giicoo/osiris/points-service/internal/delivery/restapi"
	"github.com/giicoo/osiris/points-service/internal/repository/postgres"
	"github.com/giicoo/osiris/points-service/internal/server"
	"github.com/giicoo/osiris/points-service/internal/services"
	"github.com/giicoo/osiris/points-service/pkg/logging"
	"github.com/sirupsen/logrus"
)

func RunApp() {
	logging.SetupLogging("points")
	cfg := config.SetupConfig("points")

	repository := postgres.NewRepo(cfg)
	if err := repository.Connection(); err != nil {
		logrus.Errorf("connect db: %s", err)
	}
	logrus.Info("Repo connection")
	services := services.NewServices(cfg, repository)
	controller := restapi.NewController(cfg, services)

	r := restapi.SetupRouter(controller)
	// h := cors.Default().Handler(r.Handler())

	srv := server.NewServer(cfg, r.Handler())

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

	if err := repository.CloseConnection(); err != nil {
		logrus.Errorf("close connection db: %s", err)
		return
	}
	logrus.Info("DB stop")
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
