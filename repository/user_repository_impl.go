package repository

import (
	"github.com/alifirfandi/simple-cashier-inventory/entity"
	"github.com/alifirfandi/simple-cashier-inventory/model"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	Mysql gorm.DB
}

func NewUserRepository(Mysql *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		Mysql: *Mysql,
	}
}

func (Repository UserRepositoryImpl) Profile(Request model.ProfileRequest) (Response model.ProfileResponse, Error error) {
	var user entity.User
	if Error = Repository.Mysql.Where("email = ?", Request.Email).Find(&user).Error; Error != nil {
		return Response, Error
	}
	Response.Id = user.Id
	Response.Email = user.Email
	return Response, Error
}
