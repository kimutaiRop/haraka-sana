package main

import (
	"fmt"
	"haraka-sana/config"
	oauthRoutes "haraka-sana/oauth/routes"
	ordersRoutes "haraka-sana/orders/routes"
	permissionRoutes "haraka-sana/permissions/routes"
	staffRoutes "haraka-sana/staff/routes"
	authRoutes "haraka-sana/users/routes"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func setupRouter() *gin.Engine {
	store := cookie.NewStore([]byte(os.Getenv("SECRET_KEY")))
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")
	r.Use(sessions.Sessions("my-session", store))

	r.Use(cors.New(cors.Config{
		AllowOrigins: strings.Split(os.Getenv("ALLOW_HOSTS"), ","),
	}))

	authRoutes.SessionAuth(r)
	basePath := r.Group("/api/v1")
	oauthRoutes.OauthRoutes(basePath)
	authRoutes.AuthRoutes(basePath)
	ordersRoutes.OrdersRoutes(basePath)
	staffRoutes.StaffRoutes(basePath)
	permissionRoutes.PermissionRoutes(basePath)
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
