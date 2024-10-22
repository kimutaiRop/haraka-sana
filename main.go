package main

import (
	"haraka-sana/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	models.ConnectDatabase()
	models.SeedDatabase()
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
