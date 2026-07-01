package main

import (
	"context"
	"fmt"
	"github.com/Stipulations/goAuthy/internal/config"
	"github.com/Stipulations/goAuthy/internal/database"
	"github.com/Stipulations/goAuthy/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	pool, err := database.New(ctx, cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

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

	h := handler.New(pool)
	h.RegisterRoutes(gRouter)

	fmt.Println("goAuthy is listening at:", cfg.ListenAddr)
	if err := gRouter.Run(cfg.ListenAddr); err != nil {
		panic(err)
	}
}
