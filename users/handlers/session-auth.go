package handlers

import (
	"fmt"
	"haraka-sana/config"
	"haraka-sana/helpers"
	"haraka-sana/users/models"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ShowLoginPage(c *gin.Context) {
	errorMsg := c.Query("error")

	c.HTML(http.StatusOK, "login.html", gin.H{
		"Error": errorMsg,
	})
}

func ShowSuccessPage(c *gin.Context) {
	session := sessions.Default(c)     // Get the session
	userSession := session.Get("user") // Get user ID from the session

	// Check if userSession is nil
	if userSession == nil {
		// Handle the case where the user is not found in the session
		c.Redirect(http.StatusSeeOther, "/login?error=session expired")
		return
	}

	userId := userSession.(string)   // Type assert to string
	successMsg := c.Query("message") // Get the success message from query parameters

	c.HTML(http.StatusOK, "success.html", gin.H{
		"User":    userId,
		"message": successMsg,
	})
}
func SessionLogin(c *gin.Context) {
	session := sessions.Default(c)

	// Simulate user login (you'd typically check credentials here)
	username := c.PostForm("username")
	password := c.PostForm("password")

	var loginUser models.User
	err := config.DB.Where(&models.User{Username: username}).
		Or(&models.User{Email: username}).
		First(&loginUser).
		Error

	if err != nil {
		c.Redirect(http.StatusFound, "/login?error=invalid%20credentials")
		return
	}

	if loginUser.Id == 0 {
		c.Redirect(http.StatusFound, "/login?error=invalid%20credentials")
		return
	}

	if !loginUser.Active {
		c.Redirect(http.StatusFound, "/login?error=account%20with%20username%20not%20active")
		return
	}

	if !helpers.CheckPasswordHash(password, loginUser.Password) {
		c.Redirect(http.StatusFound, "/login?error=invalid%20credentials")
		return
	}

	session.Set("user", loginUser.Email)
	fmt.Println(loginUser.Email)
	redirectURL := session.Get("redirect")
	session.Delete("redirect")
	session.Save()
	if redirectURL != nil {
		c.Redirect(http.StatusSeeOther, redirectURL.(string))
		return
	}

	c.Redirect(http.StatusSeeOther, "/success?message=login successful")

}
