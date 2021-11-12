package service

import "github.com/alifirfandi/simple-cashier-inventory/model"

type UserService interface {
	Profile(Request model.ProfileRequest) (Response model.ProfileResponse, Error error)
}
