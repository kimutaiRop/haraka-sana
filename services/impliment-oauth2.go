package services

import (
	"haraka-sana/models"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

type Grant struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GenerateAuthorizationCode(clientID, scope, redirectURI string) models.Code {
	// Generate a random string for the authorization code
	code := generateRandomString(32)
	grantCode := models.Code{}
	grantCode.Code = code
	grantCode.Scope = scope
	grantCode.RedirectURI = redirectURI
	grantCode.Expiry = time.Now().Add(10 * time.Minute)

	models.DB.Save(&grantCode)
	return grantCode
}

func CreateUniqueToken(db *gorm.DB) (*models.AuthorizationToken, error) {
	var token models.AuthorizationToken
	for {
		// Generate a new token
		authToken := generateRandomString(64)

		// Check if the token already exists
		err := db.Where("code = ?", authToken).First(&token).Error
		if err == gorm.ErrRecordNotFound {
			// Token does not exist, so we can create it
			token.Code = authToken
			token.Expiry = time.Now().Add(1 * time.Hour)

			if err := db.Save(&token).Error; err != nil {
				return nil, err
			}
			return &token, nil
		} else if err != nil {
			// Some other error occurred
			return nil, err
		}
	}

}
