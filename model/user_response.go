package model

import "time"

type UserResponse struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserListResponse struct {
	TotalData    int64          `json:"total_data"`
	TotalPage    int            `json:"total_page"`
	CurrentPage  int            `json:"current_page"`
	LimitPerPage int            `json:"limit_per_page"`
	Users        []UserResponse `json:"admins"`
}
