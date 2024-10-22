package models

import "time"

type Organization struct {
	Id              int    `json:"id"`
	ApplicationName string `json:"application_name"`
	Website         string `json:"website"`
	Logo            string `json:"logo"`
	RedirectURIs    string `json:"redirect_uris"`
	ClientId        string `json:"client_id"`
	ClientSecret    string `json:"_"`
}

type Code struct {
	Id          int
	Code        string
	Scope       string
	RedirectURI string
	Expiry      time.Time
}

type AuthorizationToken struct {
	Id     int
	Code   string
	Expiry time.Time
}
