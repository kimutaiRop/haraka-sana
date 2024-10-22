package main

import (
	"fmt"
	"haraka-sana/controller"
	"haraka-sana/models"
	"math/rand"
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
		oauth2.POST("/token", controller.AuthorizeToken)
		oauth2.POST("/client-credentials", controller.ClientCredentials)
	}
	return r
}

func main() {
	godotenv.Load()
	port := ":8080"
	rand.Seed(time.Now().UnixNano())
	mode := os.Getenv("APP_MODE")
	gin.SetMode(mode)
	models.ConnectDatabase()
	models.SeedDatabase()
	r := setupRouter()

	fmt.Println("staring at http://0.0.0.0" + port)
	r.Run(port)
}
