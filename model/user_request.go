package model

type UserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserSelectQuery struct {
	Page  int `query:"page"`
	Limit int
	Q     string `query:"q"`
}
