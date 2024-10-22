package main

import (
	"haraka-sana/controller"
	"haraka-sana/models"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		// AllowOrigins:  []string{"http://example.com"},
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:  []string{"Authorization", "X-Requested-With", "X-Request-ID", "X-HTTP-Method-Override", "Upload-Length", "Upload-Offset", "Tus-Resumable", "Upload-Metadata", "Upload-Defer-Length", "Upload-Concat", "User-Agent", "Referrer", "Origin", "Content-Type", "Content-Length"},
		ExposeHeaders: []string{"Upload-Offset", "Location", "Upload-Length", "Tus-Version", "Tus-Resumable", "Tus-Max-Size", "Tus-Extension", "Upload-Metadata", "Upload-Defer-Length", "Upload-Concat", "Location", "Upload-Offset", "Upload-Length"},
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
