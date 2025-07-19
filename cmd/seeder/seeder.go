package main

import (
	"fmt"
	"log"
	"reapp/internal/modules/user/userseeder"
	"reapp/pkg/env"
)

func main() {
	cf := env.Load("config/config.yaml")
	db, _ := env.DialMysql(cf)

	if err := userseeder.Run(db); err != nil {
		log.Fatalf("Seed failed: %v", err)
	}

	fmt.Printf("Seeder success")
}
