package middleware

import (
	"haraka-sana/config"
	oauthModels "haraka-sana/oauth/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func OrganizationTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.Request.Header.Get("Authorization")
		tokens := strings.Split(accessToken, " ")
		bearer := tokens[1]

		var organization oauthModels.OraganizationApplication
		var oauth oauthModels.AuthorizationToken

		config.DB.Where(&oauthModels.AuthorizationToken{
			Code: bearer,
		}).First(&oauth)

		if oauth.Id == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorizes App access",
			})
			c.Abort()
			return
		}

		config.DB.Where(&oauthModels.OraganizationApplication{
			Id: oauth.OrganizationAppID,
		}).First(&organization)

		if organization.Id == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorizes App access",
			})
			c.Abort()
			return
		}

		c.Set("orgazation_app", organization)
		c.Next() // Continue to the next middleware or handler if no errors

	}
}
