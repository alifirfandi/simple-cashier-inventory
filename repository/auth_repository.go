package repository

import "github.com/alifirfandi/simple-cashier-inventory/model"

type AuthRepository interface {
	Login(Request model.AuthRequest) (Response model.AuthResponse, UserExists bool, Error error)
}
