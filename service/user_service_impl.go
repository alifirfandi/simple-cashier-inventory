package service

import (
	"github.com/alifirfandi/simple-cashier-inventory/helper"
	"github.com/alifirfandi/simple-cashier-inventory/model"
	"github.com/alifirfandi/simple-cashier-inventory/repository"
	"github.com/alifirfandi/simple-cashier-inventory/validation"
	"os"
	"strconv"
)

type UserServiceImpl struct {
	UserRepository repository.UserRepository
}

func (Service UserServiceImpl) InsertUser(Request model.UserRequest) (Response model.UserResponse, Error error) {
	if Error = validation.InsertUserValidation(Request); Error != nil {
		return Response, Error
	}

	Hashed, Error := helper.GenerateHash(Request.Password)
	if Error != nil {
		return Response, Error
	}
	Request.Password = Hashed

	Response, Error = Service.UserRepository.InsertUser(Request)
	return Response, Error
}

func (Service UserServiceImpl) GetAllUser(QueryParams model.UserSelectQuery) (Response model.UserListResponse, Error error) {

	if QueryParams.Page <= 0 {
		QueryParams.Page = 1
	}

	users, TotalData, Error := Service.UserRepository.GetAllUser(QueryParams)

	limirPerPage, _ := strconv.Atoi(os.Getenv("LIMIT_PAGE"))

	Response = model.UserListResponse{
		TotalData:    TotalData,
		TotalPage:    (int(TotalData) / limirPerPage) + 1,
		LimitPerPage: limirPerPage,
		CurrentPage:  QueryParams.Page,
		Users:        users,
	}

	return Response, Error
}

func (Service UserServiceImpl) GetUserById(Id int64) (Response model.UserResponse, Error error) {

	Response, Error = Service.UserRepository.GetUserById(Id)
	return Response, Error
}

func (Service UserServiceImpl) UpdateUser(Id int64, Request model.UserRequest) (Response model.UserResponse, Error error) {
	if Error = validation.UpdateUserValidation(Request); Error != nil {
		return Response, Error
	}

	Response, Error = Service.UserRepository.UpdateUser(Id, Request)
	return Response, Error
}

func (Service UserServiceImpl) DeleteUser(Id int64) (Error error) {
	if Error = Service.UserRepository.DeleteUser(Id); Error != nil {
		return Error
	}

	return nil
}

func NewUserService(UserRepository *repository.UserRepository) UserService {
	return &UserServiceImpl{
		UserRepository: *UserRepository,
	}
}
