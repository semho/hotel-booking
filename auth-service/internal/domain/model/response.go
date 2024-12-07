package model

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"` // т.к. может не возвращаться в ответе
	User         *User  `json:"user"`
}
