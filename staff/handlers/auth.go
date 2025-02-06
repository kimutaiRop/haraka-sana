package handlers

import (
	"fmt"
	"haraka-sana/config"
	"haraka-sana/helpers"
	"haraka-sana/staff/models"
	"haraka-sana/staff/objects"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func CreateStaff(createStaff *objects.CreateStaff) (*models.Staff, gin.H) {
	if createStaff.Email == "" {
		error := gin.H{
			"field": "email",
			"error": "Email is required",
		}
		return nil, error

	}

	if createStaff.FistName == "" {
		error := gin.H{
			"field": "firstname",
			"error": "firstname is required",
		}
		return nil, error

	}

	if createStaff.LastName == "" {
		error := gin.H{
			"field": "lastname",
			"error": "firstname is required",
		}
		return nil, error

	}

	var foundStaff models.Staff
	config.DB.Where(&models.Staff{Email: createStaff.Email}).
		First(&foundStaff)
	if foundStaff.Id != 0 && foundStaff.Email == createStaff.Email {
		return nil, gin.H{"error": "Staff with email already exists"}

	}

	staff := models.Staff{
		Email:      createStaff.Email,
		FirstName:  createStaff.FistName,
		LastName:   createStaff.LastName,
		Phone:      createStaff.Phone,
		Country:    createStaff.Country,
		City:       createStaff.City,
		PositionID: createStaff.PositionID,
		Active:     false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	email_token, err := helpers.GenerateVerifyEmailToken(helpers.VerifyClaims{
		Email:       createStaff.Email,
		AccountType: "staff",
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	})

	if err != nil {
		return nil, gin.H{
			"sucess":  "account created successfully",
			"warning": "error generating email verification token"}

	}

	config.DB.Create(&staff)
	url := os.Getenv("STAFF_PORTAL") + "/auth/set-password/" + email_token
	templateData := struct {
		Name    string
		Link    string
		Company string
	}{
		Name:    foundStaff.FirstName,
		Link:    url,
		Company: os.Getenv("COMPANY_NAME"),
	}

	r := helpers.NewEmailRequest([]string{staff.Email}, "Hello "+staff.FirstName,
		"Activate your Account with Using link: "+url)
	if err := r.ParseTemplate("templates/emails/set-password.html", templateData); err != nil {
		return nil, gin.H{
			"warning": "error sending email",
			"success": "account created successfully",
		}
	}
	sent, err := r.SendEmail()
	if err != nil || !sent {
		return nil, gin.H{
			"warning": "error sending email",
			"success": "account created successfully",
		}
	}

	return &staff, nil
}

func CreateAdmin(c *gin.Context) {
	var createStaff *objects.CreateStaff

	// np other staff must exists
	var existStaff models.Staff

	config.DB.First(&existStaff)

	if existStaff.Id != 0 {
		c.JSON(500, gin.H{
			"error": "Staff Already exist cannot create new",
		})
		return
	}

	if err := c.ShouldBindJSON(&createStaff); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	staff, res := CreateStaff(createStaff)
	if res != nil {
		c.JSON(500, res)
		return
	}

	c.JSON(200, gin.H{
		"staff": staff,
	})

}

func CreateNewStaff(c *gin.Context) {
	var createStaff *objects.CreateStaff

	if err := c.ShouldBindJSON(&createStaff); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	staff, res := CreateStaff(createStaff)
	if res != nil {
		c.JSON(500, res)
		return
	}

	c.JSON(200, gin.H{
		"staff":   staff,
		"success": "account created successfully, email to set password sent",
	})

}

func StaffLogin(c *gin.Context) {

	var loginStaff objects.StaffLogin

	if err := c.ShouldBindJSON(&loginStaff); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var user models.Staff
	err := config.DB.Where(&models.Staff{Email: loginStaff.Email}).
		First(&user).
		Error

	if err != nil {
		c.JSON(400, gin.H{"error": "invalid credentials"})
		return
	}

	if user.Id == 0 {
		c.JSON(400, gin.H{"error": "invalid credentials"})
		return
	}

	if !user.Active {
		c.JSON(400, gin.H{"error": "account with username not active"})
		return
	}

	if !helpers.CheckPasswordHash(loginStaff.Password, user.Password) {
		c.JSON(400, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := helpers.GenerateToken(helpers.AuthClaims{
		ID:          user.Id,
		AccountType: "staff",
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"warning": "error logging in, please try again or contact support"})
		return
	}
	c.JSON(200, gin.H{
		"token": token,
		"user":  user,
	})
}

func SetPassword(c *gin.Context) {
	var verifyAccount *objects.StaffSetPassword

	err := c.ShouldBindJSON(&verifyAccount)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	claims, err := helpers.ValidateVerifyEmailToken(verifyAccount.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access: invalid token"})
		return
	}
	fmt.Println("claims.AccountType", claims.AccountType)
	if claims.AccountType != "staff" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access: invalid account type"})
		return
	}
	var setPassStaff models.Staff
	err = config.DB.Where(&models.Staff{Email: claims.Email}).
		Find(&setPassStaff).Error

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Could not find staff"})
		return
	}

	if setPassStaff.Id == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Could not find staff"})
		return
	}

	setPassStaff.Password = helpers.HashAndSalt([]byte(verifyAccount.Password))
	setPassStaff.VerifiedAt = time.Now()
	setPassStaff.UpdatedAt = time.Now()
	setPassStaff.Active = true

	config.DB.Save(&setPassStaff)

	c.JSON(200, gin.H{"success": "password set successfully"})
}

func UpdateStaffActiveStatus(c *gin.Context) {

	updateActiveStatus := objects.UpdateStaffActive{}

	if err := c.ShouldBindJSON(&updateActiveStatus); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var staff models.Staff

	if err := config.DB.Where("id = ?", updateActiveStatus.Id).First(&staff).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	staff.Active = updateActiveStatus.Active

	config.DB.Save(&staff)

	c.JSON(200, gin.H{"success": "staff active status updated"})
}

func StaffRequestPasswordReset(c *gin.Context) {

	var input objects.StaffRequestPasswordReset
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var staff models.Staff

	getErr := config.DB.Where(&models.Staff{
		Id: input.Id,
	}).First(&staff).Error

	if getErr != nil {
		c.JSON(400, gin.H{"error": "staff not found"})
		return
	}

	config.DB.Model(&models.Staff{}).Where("id = ?", input.Id).Updates(
		models.Staff{
			Active: false,
		},
	)

	email_token, err := helpers.GenerateVerifyEmailToken(helpers.VerifyClaims{
		Email:       staff.Email,
		AccountType: "staff",
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	})

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error":  "error generating email reset token",
			"status": "failed"})
		return
	}

	templateData := struct {
		Name string
		URL  string
	}{
		Name: staff.FirstName,
		URL:  os.Getenv("STAFF_PORTAL") + "/auth/set-password/" + email_token,
	}

	r := helpers.NewEmailRequest([]string{staff.Email}, "Hello "+staff.FirstName,
		"Set your password by clicking on this link"+

			" "+os.Getenv("STAFF_PORTAL")+"/auth/set-password/"+email_token)

	if err := r.ParseTemplate("templates/set-password.html", templateData); err == nil {
		r.SendEmail()
	}

	c.JSON(200, gin.H{"success": "email to set password sent"})
}
