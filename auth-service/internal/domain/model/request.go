package model

type RegisterRequest struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
	Phone     *string
}

type LoginRequest struct {
	Email    string
	Password string
}
