package middleware

import (
	"fmt"
	"haraka-sana/config"
	"haraka-sana/helpers"
	"haraka-sana/users/models"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userSession := session.Get("user")

		if userSession == nil {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		fmt.Println(userSession)

		userEmail := userSession.(string)

		var user models.User

		getErr := config.DB.Where(&models.User{Email: userEmail}).First(&user).Error

		if getErr != nil {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

func JWTAuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		accessToken := c.Request.Header.Get("Authorization")
		token, err := helpers.ValidateToken(accessToken)
		if err != nil {
			fmt.Print(err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Auth token is invalid",
			})
			c.Abort()
			return
		}
		if token.AccountType != "user" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Account type",
			})
			c.Abort()
			return
		}

		var user models.User

		getErr := config.DB.Where(&models.User{Id: token.ID}).First(&user).Error

		if getErr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Auth token is invalid",
			})
			c.Abort()
			return
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
