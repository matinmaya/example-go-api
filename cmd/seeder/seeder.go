package main

import (
	"fmt"
	"log"
	"reapp/config"
	"reapp/internal/modules/user/userseeder"
)

func main() {
	cf := config.Load("config/config.yaml")
	db, _ := config.DialMysql(cf)

	if err := userseeder.Run(db); err != nil {
		log.Fatalf("Seed failed: %v", err)
	}

	fmt.Printf("Seeder success")
}
