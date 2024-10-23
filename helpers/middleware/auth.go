package middleware

import (
	"haraka-sana/config"
	"haraka-sana/helpers"
	"haraka-sana/users/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		accessToken := c.Request.Header.Get("Authorization")
		token, err := helpers.ValidateToken(accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Auth token is invalid",
			})
		}
		if token.AccountType != "user" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Account type",
			})
		}

		var user models.User

		getErr := config.DB.Where(&models.User{Id: token.ID}).First(&user).Error

		if getErr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Auth token is invalid",
			})
		}
		if !user.Active {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Account is not active",
			})
		}
		c.Set("user", user)
		c.Next()

	}
}
