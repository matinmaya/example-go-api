package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm/logger"
)

type ErrorLog struct {
	LogFile *os.File
	logger.Interface
}

func NewErrorLog() logger.Interface {
	today := time.Now().Format("2006-01-02")
	logDir := "logs/sql"
	_ = os.MkdirAll(logDir, 0755)

	filePath := filepath.Join(logDir, fmt.Sprintf("error-logger-%s.log", today))
	logFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to open log file: %v", err))
	}

	newLogger := log.New(logFile, "", log.LstdFlags)

	return &ErrorLog{
		LogFile: logFile,
		Interface: logger.New(newLogger, logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		}),
	}
}
