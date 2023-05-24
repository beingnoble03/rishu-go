package main

import (
	"example/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/ok", handler.Ok)
	r.GET("/health", handler.Hand)

	r.Run()
}
