package objects

type CreateStaff struct {
	Email      string `json:"email"`
	FistName   string `json:"firstname"`
	LastName   string `json:"lastname"`
	Phone      string `json:"phone"`
	IDNumber   string `json:"id_number"`
	Country    string `json:"country"`
	City       string `json:"city"`
	PositionID int    `json:"position_id"`
}
type StaffSetPassword struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type StaffLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type StaffRequestPasswordReset struct {
	Id int `json:"id"`
}

type UpdateStaffActive struct {
	Active bool `json:"active"`
	Id     int  `json:"id"`
}
