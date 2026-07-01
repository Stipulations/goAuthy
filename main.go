package main

import (
	"context"
	"fmt"

	"github.com/Stipulations/goAuthy/internal/config"
	"github.com/Stipulations/goAuthy/internal/database"
	"github.com/Stipulations/goAuthy/internal/handler"
	"github.com/Stipulations/goAuthy/internal/jwt"
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
	if err := database.Migrate(ctx, pool); err != nil {
		panic(err)
	}
	if err := database.CreateAdmin(ctx, pool); err != nil {
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

	tokenSvc := jwt.New(cfg.JWTSecret)
	h := handler.New(pool, tokenSvc)
	h.RegisterRoutes(gRouter)

	fmt.Println("goAuthy is listening at:", cfg.ListenAddr)
	if err := gRouter.Run(cfg.ListenAddr); err != nil {
		panic(err)
	}
}
