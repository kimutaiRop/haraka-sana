package main

import (
	"fmt"
	"haraka-sana/config"
	oauthRoutes "haraka-sana/oauth/routes"
	authRoutes "haraka-sana/users/routes"
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
	oauthRoutes.OauthRoutes(basePath)
	authRoutes.AuthRoutes(basePath)
	return r
}

func main() {
	godotenv.Load()
	port := ":8080"
	rand.Seed(time.Now().UnixNano())
	mode := os.Getenv("APP_MODE")
	gin.SetMode(mode)
	config.ConnectDatabase()
	config.SeedDatabase()
	r := setupRouter()

	fmt.Println("staring at http://0.0.0.0" + port)
	r.Run(port)
}
