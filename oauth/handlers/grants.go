package handlers

import (
	"haraka-sana/config"
	"haraka-sana/helpers"
	"haraka-sana/oauth/models"
	"haraka-sana/oauth/objects"
	"haraka-sana/oauth/services"
	authModel "haraka-sana/users/models"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func errorPage() string {
	return "/auth/error"
}

func AuthorizeCode(c *gin.Context) {
	query := c.Request.URL.Query()
	userContext, _ := c.Get("user")

	user := userContext.(authModel.User)
	redirect_uri := query.Get("redirect_uri")
	if redirect_uri == "" {
		c.Redirect(http.StatusFound, errorPage()+"?error=redirect_uri not found")
		return
	}
	grant_type := query.Get("grant_type")
	if grant_type == "" {
		c.Redirect(http.StatusFound, errorPage()+"?error=grant_type not found")
		return
	}
	client_id := query.Get("client_id")
	if client_id == "" {
		c.Redirect(http.StatusFound, errorPage()+"?error=client_id not found")
		return
	}
	organization := models.OraganizationApplication{}
	config.DB.Where(&models.OraganizationApplication{
		ClientId: client_id,
	}).First(&organization)

	if organization.Id == 0 {
		c.Redirect(http.StatusFound, errorPage()+"?error=Organization with client_id not found")
		return
	}
	if grant_type == "code" {
		scope := query.Get("scope")
		redirects := strings.Split(organization.RedirectURIs, ",")
		if !helpers.Contains(redirects, redirect_uri) {
			c.Redirect(http.StatusFound, errorPage()+"?error=Invalid redirect_uri")
			return
		}
		code := services.GenerateAuthorizationCode(client_id, scope, redirect_uri)
		redirect_uri = redirect_uri + "?code=" + code.Code

		code.OrganizationAppID = organization.Id
		code.UserId = user.Id

		config.DB.Create(&code)

		c.Redirect(http.StatusFound, redirect_uri)
	}
	c.Redirect(http.StatusFound, errorPage()+"?error=Invalid grant_type")
}

func AuthorizeToken(c *gin.Context) {

	var tokenAuth objects.TokenAuth

	err := c.ShouldBindJSON(&tokenAuth)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	redirect_uri := tokenAuth.RedirectURI

	grant_type := tokenAuth.GrantType
	client_id := tokenAuth.ClientID
	organization := models.OraganizationApplication{}
	config.DB.Where(&models.OraganizationApplication{
		ClientId: client_id,
	}).First(&organization)

	if organization.Id == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization with client_id not found"})
		return
	}
	if grant_type == "authorization_code" {
		code := tokenAuth.Code
		scope := tokenAuth.Scope
		redirects := strings.Split(organization.RedirectURIs, ",")
		if !helpers.Contains(redirects, redirect_uri) {
			c.JSON(http.StatusForbidden, gin.H{
				"errror": "Invalid redirect_uri",
			})
			return
		}
		var authCode *models.Code
		config.DB.Where(&models.Code{
			Code: code,
		}).First(&authCode)
		if authCode.Id == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"errror": "Invalid code or expired",
			})
			return
		}
		if authCode.Scope != scope {
			c.JSON(http.StatusForbidden, gin.H{
				"errror": "Invalid scope",
			})
			return
		}
		if redirect_uri != authCode.RedirectURI {
			c.JSON(http.StatusForbidden, gin.H{
				"errror": "Invalid redirect_uri",
			})
			return
		}
		if authCode.Id == 0 || authCode.Expiry.Before(time.Now()) {
			c.JSON(http.StatusForbidden, gin.H{
				"errror": "Invalid code or Expired",
			})
			return
		}

		token, err := services.CreateUniqueToken(config.DB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errror": "Error in our end",
			})
			return
		}
		token.OrganizationAppID = organization.Id
		config.DB.Create(&token)
		if err := config.DB.Delete(&authCode).Error; err != nil {
			// Handle deletion errors
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Failed to delete code"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"access_token": token.Code,
			"expires_in":   60 * 60,
		})
		return
	}
	c.JSON(http.StatusForbidden, gin.H{
		"errror": "invalid grant_type",
	})
}

func ClientCredentials(c *gin.Context) {

	var clientCred *objects.ClientCred

	err := c.ShouldBindJSON(&clientCred)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if clientCred.GrantType != "authorization_code" {
		c.JSON(http.StatusForbidden, gin.H{
			"errror": "invalid grant_type",
		})
	}
	organization := models.OraganizationApplication{}
	config.DB.Where(&models.OraganizationApplication{
		ClientId: clientCred.ClientId,
	}).First(&organization)

	if organization.Id == 0 || organization.ClientSecret != clientCred.ClientSecret {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Organization not found",
		})
		return
	}
	token, err := services.CreateUniqueToken(config.DB)
	token.OrganizationAppID = organization.Id

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errror": "Error in our end",
		})
		return
	}
	config.DB.Create(&token)

	c.JSON(http.StatusOK, gin.H{
		"access_token": token.Code,
		"expires_in":   60 * 60,
	})
	return
}
