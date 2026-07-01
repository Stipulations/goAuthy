package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	skidlistener := "127.0.0.1:6969"

	gin.SetMode(gin.ReleaseMode)
	gRouter := gin.Default()

	gRouter.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"hello": "world",
		})
	})

	fmt.Println("goAuthy is listening at:", skidlistener)
	gRouter.Run(skidlistener)
}