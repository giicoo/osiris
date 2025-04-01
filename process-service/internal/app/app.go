package app

import (
	"os"
	"os/signal"
	"time"

	"github.com/giicoo/osiris/process-service/internal/config"
	"github.com/giicoo/osiris/process-service/internal/infrastructure/rabbitmq"
	"github.com/giicoo/osiris/process-service/internal/repository/postgres"
	"github.com/giicoo/osiris/process-service/internal/services"
	"github.com/giicoo/osiris/process-service/pkg/logging"
	"github.com/sirupsen/logrus"
)

func RunApp() {
	logging.SetupLogging("process")
	cfg := config.SetupConfig("process")

	rabmq := rabbitmq.NewAlertProducing(cfg)
	if err := rabmq.InitAlertProducing(); err != nil {
		logrus.Fatal("init rabbitmq: ", err)
	}
	repository := postgres.NewRepo(cfg)
	if err := repository.Connection(); err != nil {
		logrus.Fatal("connect db: %w", err)
	}
	logrus.Info("Repo connection")
	services := services.NewServices(cfg, repository, rabmq)

	go func() {
		msgs, err := rabmq.ConsumeMessage()
		if err != nil {
			logrus.Fatal(err)
		}
		for msg := range msgs {
			if err := services.Processing(msg.Body); err != nil {
				logrus.Error(err)
			}
			time.Sleep(1 * time.Second)
			if err := services.Processing(msg.Body); err != nil {
				logrus.Error(err)
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
}
