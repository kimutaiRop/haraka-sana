package objects

type CreateUser struct {
	Email           string `json:"email"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confim_password"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RequestPasswordReset struct {
	Email string `json:"email"`
}
