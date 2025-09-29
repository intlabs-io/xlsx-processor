package main

import (
	"net/http"

	"xlsx-processor/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	err := initSentry(router)
	if err != nil {
		panic(err)
	}

	router.POST("/xlsx-processor/paginate", routes.Paginate)

	router.POST("/xlsx-processor/transform", routes.Transform)

	router.POST("/xlsx-processor/transform-json", routes.TransformJson)

	router.GET("/xlsx-processor/healthz/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	router.GET("/xlsx-processor/healthz/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	router.GET("/xlsx-processor/healthz/exception", func(c *gin.Context) {
		panic("exception")
	})

	router.Run(":8080")
}
