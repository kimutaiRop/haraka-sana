package middleware

import (
	"haraka-sana/config"
	"haraka-sana/users/models"
	"haraka-sana/users/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		accessToken := c.Request.Header.Get("Authorization")
		token, err := services.ValidateToken(accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Auth token is invalid",
			})
		}
		var user models.User

		getErr := config.DB.Where("email = ?", token.Email).First(&user).Error

		if getErr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Auth token is invalid",
			})
		}
		c.Set("user", user)
		c.Next()

	}
}
