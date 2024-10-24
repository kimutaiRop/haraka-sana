package objects

type CreateUser struct {
	Email           string `json:"email"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type VerifyAccount struct {
	Token string `json:"token"`
}

type RequestPasswordReset struct {
	Email string `json:"email"`
}

type SetPassword struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}
