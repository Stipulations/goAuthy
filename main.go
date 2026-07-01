package main

import (
	"fmt"
	"net/http"
	"github.com/Stipulations/goAuthy/internal"
	"github.com/gin-gonic/gin"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		panic(err)
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