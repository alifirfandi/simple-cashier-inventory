package service

import "github.com/alifirfandi/simple-cashier-inventory/model"

type AuthService interface {
	Login(Request model.AuthRequest) (Response model.AuthResponse, UserExists bool, Error error)
}
