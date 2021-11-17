package service

import (
	"errors"
	"fmt"
	"github.com/alifirfandi/simple-cashier-inventory/helper"
	"github.com/alifirfandi/simple-cashier-inventory/model"
	"github.com/alifirfandi/simple-cashier-inventory/repository"
	"github.com/alifirfandi/simple-cashier-inventory/validation"
	"math"
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

	userExist := Service.UserRepository.CheckUserExist(Request.Email)
	if userExist {
		return Response, errors.New("EMAIL_REGISTERED")
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

	QueryParams.Q = fmt.Sprintf("%%%s%%", QueryParams.Q)

	limitPerPage, Error := strconv.Atoi(os.Getenv("LIMIT_PER_PAGE"))
	currentPage := QueryParams.Page
	QueryParams.Page = (QueryParams.Page - 1) * limitPerPage

	users, TotalData, Error := Service.UserRepository.GetAllUser(QueryParams)

	Response = model.UserListResponse{
		TotalData:    TotalData,
		TotalPage:    int(math.Ceil(float64(TotalData) / float64(limitPerPage))),
		LimitPerPage: limitPerPage,
		CurrentPage:  currentPage,
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

	userExist := Service.UserRepository.CheckUserExist(Request.Email)
	if userExist {
		return Response, errors.New("EMAIL_REGISTERED")
	}

	Hashed, Error := helper.GenerateHash(Request.Password)
	Request.Password = Hashed

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
