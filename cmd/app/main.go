package main

import (
	"fmt"
	"reapp/internal/helpers/redishelper"
	"reapp/internal/provider"
	"reapp/pkg/env"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	cf := env.Load("config/config.yaml")
	db, _ := env.DialMysql(cf)
	redisClient, _ := env.DialRedis(cf)
	redishelper.InitRedis(redisClient)

	provider.NewProvider(r, db, cf).RegisterServiceProvider().RegisterRouteProvider()

	r.Run(fmt.Sprintf(":%d", cf.App.Port))
}
