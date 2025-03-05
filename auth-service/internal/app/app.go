package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/giicoo/osiris/auth-service/internal/config"
	"github.com/giicoo/osiris/auth-service/internal/delivery/restapi"
	"github.com/giicoo/osiris/auth-service/internal/infrastructure"
	"github.com/giicoo/osiris/auth-service/internal/repository/postgres"
	"github.com/giicoo/osiris/auth-service/internal/server"
	"github.com/giicoo/osiris/auth-service/internal/services"
	"github.com/giicoo/osiris/auth-service/pkg/logging"
	"github.com/sirupsen/logrus"
)

func RunApp() {
	logging.SetupLogging("auth-service")
	cfg := config.SetupConfig("auth-service")
	logrus.Info(cfg)
	repository := postgres.NewRepo(cfg)
	if err := repository.Connection(); err != nil {
		logrus.Fatal("connect db: %w", err)
	}
	logrus.Info("Repo connection")
	yandexAPI := infrastructure.NewYandexAPI(cfg)
	services := services.NewServices(cfg, repository, yandexAPI)
	controller := restapi.NewController(cfg, services)

	r := restapi.SetupRouter(controller)

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
