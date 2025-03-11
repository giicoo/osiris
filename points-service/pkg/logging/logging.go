package logging

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func SetupLogging(service string) {
	logrus.SetFormatter(NewFormatter(false))
	fileLogger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("./logs/%s/%s.log", service, service), // Путь к файлу
		MaxSize:    10,                                                // Максимальный размер файла в МБ до ротации
		MaxBackups: 3,                                                 // Максимальное количество старых файлов
		MaxAge:     7,                                                 // Максимальное количество дней для хранения файлов
		Compress:   true,                                              // Сжатие старых логов
	}

	logrus.SetOutput(io.MultiWriter(fileLogger, os.Stdout))
	// logrus.SetOutput(fileLogger)
}
