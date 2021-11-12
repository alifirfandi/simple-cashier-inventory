package service

import (
	"github.com/alifirfandi/simple-cashier-inventory/model"
	"github.com/alifirfandi/simple-cashier-inventory/repository"
)

type UserServiceImpl struct {
	UserRepository repository.UserRepository
}

func NewUserService(UserRepository *repository.UserRepository) UserService {
	return &UserServiceImpl{
		UserRepository: *UserRepository,
	}
}

func (Service UserServiceImpl) Profile(Request model.ProfileRequest) (Response model.ProfileResponse, Error error) {
	Response, Error = Service.UserRepository.Profile(Request)
	return Response, Error
}
