package main

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net/http"
)

type Config struct {
	AppState   string `env:"AppState,required"`
	ListenAddr string `env:"ListenerAddr,required"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Please create a .env file")
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		panic("Please fill in the AppState and or ListenerAddr in the .env file")
	}

	var gRouter *gin.Engine

	switch cfg.AppState {
	case "prod":
		gin.SetMode(gin.ReleaseMode)
		gRouter = gin.New()
		gRouter.Use(gin.Recovery())
	case "dev":
		gRouter = gin.Default()
	default:
		panic("AppState has to be either 'prod' or 'dev'")
	}

	gRouter.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"hello": "world",
		})
	})

	fmt.Println("goAuthy is listening at:", cfg.ListenAddr)
	if err := gRouter.Run(cfg.ListenAddr); err != nil {
		panic(err)
	}
}
