package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
}

func SetupConfig(service string) *Config {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Не удалось прочитать конфигурационный файл: %v", err)
	}

	var config Config

	// Используем BindEnv, чтобы связать переменные окружения с соответствующими полями
	// viper.BindEnv("server.port", "SERVER_PORT")

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Ошибка при чтении конфигурации: %v", err)
	}

	return &config
}
