package model

type AuthResponse struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"-"`
	Role        string `json:"role"`
	AccessToken string `json:"access_token"`
}
