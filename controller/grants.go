package controller

import (
	"haraka-sana/helpers"
	"haraka-sana/models"
	"haraka-sana/objects"
	"haraka-sana/services"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthorizeCode(c *gin.Context) {
	query := c.Request.URL.Query()
	redirect_uri := query.Get("redirect_uri")

	grant_type := query.Get("grant_type")
	client_id := query.Get("client_id")
	organization := models.Organization{}
	models.DB.Where("id = ?", client_id).First(&organization)

	if organization.Id == 0 {
		c.JSON(http.StatusUnauthorized, redirect_uri+"?error=Organization with client_id not found")
		return
	}
	if grant_type == "code" {
		scope := query.Get("scope")
		redirects := strings.Split(organization.RedirectURIs, ",")
		if !helpers.Contains(redirects, redirect_uri) {
			c.JSON(http.StatusUnauthorized, redirect_uri+"?error=Invalid redirect_uri")
			return
		}
		code := services.GenerateAuthorizationCode(client_id, scope, redirect_uri)
		redirect_uri = redirect_uri + "?code=" + code.Code
		c.Redirect(http.StatusOK, redirect_uri)
	}
	c.Redirect(http.StatusUnauthorized, redirect_uri+"?error=Invalid grant_type")
}

func AuthorizeToken(c *gin.Context) {
	query := c.Request.URL.Query()
	redirect_uri := query.Get("redirect_uri")

	grant_type := query.Get("grant_type")
	client_id := query.Get("client_id")
	organization := models.Organization{}
	models.DB.Where("id = ?", client_id).First(&organization)

	if organization.Id == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization with client_id not found"})
		return
	}
	if grant_type == "authorization_code" {
		code := query.Get("code")
		scope := query.Get("scope")
		redirects := strings.Split(organization.RedirectURIs, ",")
		if !helpers.Contains(redirects, redirect_uri) {
			c.JSON(http.StatusForbidden, gin.H{
				"errror": "Invalid redirect_uri",
			})
			return
		}
		authCode := models.Code{}
		models.DB.Where("code = ?", code).First(&authCode)
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
		if authCode.Id == 0 || authCode.Expiry.After(time.Now()) {
			c.JSON(http.StatusForbidden, gin.H{
				"errror": "Invalid code or Expired",
			})
			return
		}

		token, err := services.CreateUniqueToken(models.DB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errror": "Error in our end",
			})
			return
		}
		models.DB.Save(&token)

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

	var clientCred objects.ClientCred

	if err := c.ShouldBindJSON(&clientCred); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if clientCred.GrantType != "authorization_code" {
		c.JSON(http.StatusForbidden, gin.H{
			"errror": "invalid grant_type",
		})
	}
	organization := models.Organization{}
	models.DB.Where("id = ?", clientCred.ClientId).First(&organization)

	if organization.Id == 0 || organization.ClientSecret != clientCred.ClientSecret {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Organization not found",
		})
		return
	}
	token, err := services.CreateUniqueToken(models.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errror": "Error in our end",
		})
		return
	}
	models.DB.Save(&token)

	c.JSON(http.StatusOK, gin.H{
		"access_token": token.Code,
		"expires_in":   60 * 60,
	})
	return
}
