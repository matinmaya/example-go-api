package main

import (
	"fmt"

	"reapp/config"
	"reapp/internal/modules/customer"
	"reapp/internal/modules/user/usermigration"
	"reapp/pkg/base/basemodel"
)

func main() {
	configPath := fmt.Sprintf("config/application/config.%s.yaml", "debug")
	cf := config.Load(configPath)
	db, _ := config.DialMysql(cf)

	db.AutoMigrate(&basemodel.TableLog{}, &basemodel.HttpLog{})
	usermigration.Migrate(db)
	db.AutoMigrate(&customer.Customer{})

	fmt.Printf("migration success\n")
}
