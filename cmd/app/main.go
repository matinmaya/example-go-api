package main

import (
	"fmt"
	"reapp/config"
	"reapp/internal/helpers/redishelper"
	"reapp/internal/provider"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	env := gin.Mode()
	configPath := fmt.Sprintf("config/application/config.%s.yaml", env)
	cf := config.Load(configPath)
	db, _ := config.DialMysql(cf)
	redisClient, _ := config.DialRedis(cf)
	redishelper.InitRedis(redisClient)

	provider.NewProvider(r, db, cf).RegisterServiceProvider().RegisterRouteProvider()

	r.Run(fmt.Sprintf(":%d", cf.App.Port))
}
