package repository

import (
	"github.com/alifirfandi/simple-cashier-inventory/entity"
	"github.com/alifirfandi/simple-cashier-inventory/model"

	"gorm.io/gorm"
)

type AuthRepositoryImpl struct {
	Mysql gorm.DB
}

func NewAuthRepository(Mysql *gorm.DB) AuthRepository {
	return &AuthRepositoryImpl{
		Mysql: *Mysql,
	}
}

func (Repository AuthRepositoryImpl) Login(Request model.AuthRequest) (Response model.AuthResponse, UserExists bool, Error error) {
	var user entity.User
	Error = Repository.Mysql.Where("email = ?", Request.Email).Find(&user).Error
	if user.Email == "" {
		return Response, false, Error
	}
	Response.Id = user.Id
	Response.Email = user.Email
	Response.Password = user.Password
	return Response, true, Error
}
