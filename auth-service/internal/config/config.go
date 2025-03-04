package config

import (
	"log"
	"strings"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"fmt"
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

	ClientID     string `mapstructure:"CLIENT_ID"`
	ClientSecret string `mapstructure:"CLIENT_SECRET"`
	RedirectURL  string `mapstructure:"redirect_url"`
}
func SetupConfig(service string) *Config {
	var config Config
	SetupConfigFile(&config, service)
	SetupConfigENV(&config, service)
	return &config

}
func SetupConfigFile(config *Config, service string){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()
	

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Не удалось прочитать конфигурационный файл: %v", err)
	}

	

	// Используем BindEnv, чтобы связать переменные окружения с соответствующими полями
	viper.BindEnv("db.user", "DB_USER")
	viper.BindEnv("db.password", "DB_PASSWORD")
	viper.BindEnv("db.db", "DB_DB")

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Ошибка при чтении конфигурации: %v", err)
	}

}

func SetupConfigENV(config *Config, service string) {
	_ = godotenv.Load()

	viper.SetConfigFile(".env") // Явно указываем, что используем .env
	viper.SetConfigType("env")  // Тип файла
	viper.AutomaticEnv()        // Загружаем переменные окружения

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	
	if err := viper.Unmarshal(config); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	fmt.Println("Client ID:", config.ClientID)
	fmt.Println("Client Secret:", config.ClientSecret)
	fmt.Println("Redirect URL:", config.RedirectURL)
}