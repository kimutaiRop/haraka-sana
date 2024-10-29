package models

import (
	authModel "haraka-sana/users/models"
	"time"
)

type OrganizationApplication struct {
	Id              int    `gorm:"primary_key" json:"id"`
	ApplicationName string `json:"application_name"`
	Website         string `json:"website"`
	Logo            string `json:"logo"`
	RedirectURIs    string `json:"redirect_uris"`
	ClientId        string `json:"client_id"`
	ClientSecret    string `json:"-"`

	UserId int             `json:"-"`
	User   *authModel.User `json:"-" gorm:"foreignKey:UserId"`

	EventsCallbackUrl string `json:"events_callback_url"`
	Approved          bool   `json:"approved"`
}

type Code struct {
	Id                int `gorm:"primary_key"`
	Code              string
	Scope             string
	RedirectURI       string
	Expiry            time.Time
	OrganizationAppID int
	UserId            int
	OrganizationApp   *OrganizationApplication `gorm:"foreignKey:OrganizationAppID"`
	User              *authModel.User          `gorm:"foreignKey:UserId"`
}

type AuthorizationToken struct {
	Id                int `gorm:"primary_key"`
	Code              string
	Expiry            time.Time
	OrganizationAppID int
	OrganizationApp   *OrganizationApplication `gorm:"foreignKey:OrganizationAppID"`
}
