package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/giicoo/osiris/notification-service/internal/config"
	"github.com/giicoo/osiris/notification-service/internal/delivery/restapi"
	"github.com/giicoo/osiris/notification-service/internal/infrastructure/rabbitmq"
	"github.com/giicoo/osiris/notification-service/internal/server"
	"github.com/giicoo/osiris/notification-service/internal/services"
	"github.com/giicoo/osiris/notification-service/pkg/logging"
	"github.com/sirupsen/logrus"
)

func RunApp() {
	logging.SetupLogging("notifications")
	cfg := config.SetupConfig("notifications")

	rabmq := rabbitmq.NewAlertProducing(cfg)
	if err := rabmq.InitAlertProducing(); err != nil {
		logrus.Fatal("init rabbitmq: ", err)
	}
	services := services.NewServices(cfg, rabmq)
	controller := restapi.NewController(cfg, services)

	r := restapi.SetupRouter(controller)

	srv := server.NewServer(cfg, r.Handler())

	go func() {
		msgs, err := rabmq.ConsumeMessage()
		if err != nil {
			logrus.Fatal(err)
		}
		for msg := range msgs {
			if err := services.Processing(msg.Body); err != nil {
				logrus.Error(err)
			}
		}
	}()
	logrus.Info("Rabbit Start")

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
	if err := rabmq.Close(); err != nil {
		logrus.Error(err)
	}

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
