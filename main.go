package main

import (
	"haraka-sana/controller"
	"haraka-sana/models"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: strings.Split(os.Getenv("ALLOW_HOSTS"), ","),
	}))

	basePath := r.Group("/api/v1")

	{
		oauth2 := basePath.Group("/oauth2")
		oauth2.GET("/authorize", controller.AuthorizeCode)
		oauth2.GET("/token", controller.AuthorizeToken)
		oauth2.GET("/client-credentials", controller.ClientCredentials)
	}
	return r
}

func main() {
	godotenv.Load()

	rand.Seed(time.Now().UnixNano())
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
