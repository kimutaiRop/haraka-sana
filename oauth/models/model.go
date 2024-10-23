package models

import (
	authModel "haraka-sana/users/models"
	"time"
)

type OraganizationApplication struct {
	Id              int    `gorm:"primary_key" json:"id"`
	ApplicationName string `json:"application_name"`
	Website         string `json:"website"`
	Logo            string `json:"logo"`
	RedirectURIs    string `json:"redirect_uris"`
	ClientId        string `json:"client_id"`
	ClientSecret    string `json:"-"`

	UserId int            `json:"-"`
	User   authModel.User `json:"user" gorm:"foreignKey:UserId"`
}

type Code struct {
	Id          int `gorm:"primary_key"`
	Code        string
	Scope       string
	RedirectURI string
	Expiry      time.Time
}

type AuthorizationToken struct {
	Id     int `gorm:"primary_key"`
	Code   string
	Expiry time.Time
}
