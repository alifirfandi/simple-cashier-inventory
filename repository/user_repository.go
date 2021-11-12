package repository

import "github.com/alifirfandi/simple-cashier-inventory/model"

type UserRepository interface {
	Profile(Request model.ProfileRequest) (Response model.ProfileResponse, Error error)
}
