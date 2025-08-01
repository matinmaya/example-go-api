package main

import (
	"fmt"
	"log"
	"reapp/config"
	"reapp/internal/modules/user/userseeder"
)

func main() {
	configPath := fmt.Sprintf("config/application/config.%s.yaml", "debug")
	cf := config.Load(configPath)
	db, _ := config.DialMysql(cf)

	if err := userseeder.Run(db); err != nil {
		log.Fatalf("Seed failed: %v", err)
	}

	fmt.Printf("Seeder success")
}
