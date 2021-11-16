package service

import "github.com/alifirfandi/simple-cashier-inventory/model"

type UserService interface {
	InsertUser(Request model.UserRequest) (Response model.UserResponse, Error error)
	GetAllUser(QueryParams model.UserSelectQuery) (Response model.UserListResponse, Error error)
	GetUserById(Id int64) (Response model.UserResponse, Error error)
	UpdateUser(Id int64, Request model.UserRequest) (Response model.UserResponse, Error error)
	DeleteUser(Id int64) (Error error)
}
