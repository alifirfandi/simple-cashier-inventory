package model

type AuthResponse struct {
	Id          int64  `json:"-"`
	Email       string `json:"email"`
	Password    string `json:"-"`
	AccessToken string `json:"access_token"`
}
