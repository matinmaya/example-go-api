package main

import (
	"fmt"
	"reapp/config"
	"reapp/internal/provider"
	"reapp/pkg/helpers/loghelper"
	"reapp/pkg/helpers/redishelper"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	env := gin.Mode()
	configPath := fmt.Sprintf("config/application/config.%s.yaml", env)
	cf := config.Load(configPath)
	loghelper.InitLogger(cf.Log.Filename)
	db, _ := config.DialMysql(cf)
	redisClient, _ := config.DialRedis(cf)
	redishelper.InitRedis(redisClient, cf.Redis.RepoCacheTTL)

	provider.NewProvider(r, db, cf).RegisterServiceProvider().RegisterRouteProvider()

	r.Run(fmt.Sprintf(":%d", cf.App.Port))
}
