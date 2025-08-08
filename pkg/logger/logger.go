package logger

import (
	"log"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

func InitLogger(linkname string) {
	writer, err := rotatelogs.New(
		"logs/app/%Y-%m-%d.log",
		rotatelogs.WithLinkName(linkname),         // symlink to latest log
		rotatelogs.WithMaxAge(7*24*time.Hour),     // keep logs 7 days
		rotatelogs.WithRotationTime(24*time.Hour), // rotate daily
	)
	if err != nil {
		log.Fatalf("Failed to initialize log file rotator: %v", err)
	}

	log.SetOutput(writer)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
