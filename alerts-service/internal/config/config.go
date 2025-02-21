package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Db       string `mapstructure:"db"`
	} `mapstructure:"db"`

	Server struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`

	Rabbitmq struct {
		Exchange struct {
			Name       string `mapstructure:"name"`
			Type       string `mapstructure:"type"`
			Durability bool   `mapstructure:"durability"`
		} `mapstructure:"exchange"`

		Queue struct {
			Name       string `mapstructure:"name"`
			Durability bool   `mapstructure:"durability"`
		}
	} `mapstructure:"rabbitmq"`
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
	viper.BindEnv("db.user", "DB_USER")
	viper.BindEnv("db.password", "DB_PASSWORD")
	viper.BindEnv("db.db", "DB_DB")

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Ошибка при чтении конфигурации: %v", err)
	}

	return &config
}
