package rest

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/spf13/viper"
	"time"
)

var (
	restServer *gin.Engine
)

func Init() error {
	binding.Validator = new(defaultValidator)
	restServer = gin.Default()

	// CORS
	restServer.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{},
		AllowCredentials: true,
		AllowOriginFunc:  nil,
		MaxAge:           12 * time.Hour,
	}))

	return nil
}

func RunServer() {
	host := viper.GetString("rest.host")
	port := viper.GetString("rest.port")
	restServer.Run(host + ":" + port)
}
