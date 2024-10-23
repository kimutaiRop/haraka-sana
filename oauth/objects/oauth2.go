package objects

type TokenAuth struct {
	RedirectURI string `json:"redirect_uri"`
	GrantType   string `json:"grant_type"`
	ClientID    string `json:"client_id"`
	Code        string `json:"code"`
	Scope       string `json:"scope"`
}

type ClientCred struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

type CreateApp struct {
	ApplicationName string `json:"application_name"`
	Website         string `json:"website"`
	Logo            string `json:"logo"`
	RedirectURIs    string `json:"redirect_uris"`
}
