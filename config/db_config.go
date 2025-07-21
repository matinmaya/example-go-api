package config

import (
	"fmt"
	"reapp/internal/logger"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DialMysql(cf *Config) (*gorm.DB, error) {

	db, err := gorm.Open(mysql.Open(cf.GetDsn()), &gorm.Config{
		Logger: logger.NewGormLogger(),
	})
	if err != nil {
		return nil, err
	}

	PingDB(db, cf.GetDsn())

	return db, nil
}

func PingDB(db *gorm.DB, dsn string) {
	sql, err := db.DB()
	if err != nil {
		panic(err)
	}
	if err := sql.Ping(); err != nil {
		fmt.Println("Error connecting to database:", err)
	}

	dsnParts := strings.Split(dsn, "/")
	dbName := strings.Split(dsnParts[len(dsnParts)-1], "?")[0]

	fmt.Printf("Connected to MySQL database: %s\n", dbName)
}
