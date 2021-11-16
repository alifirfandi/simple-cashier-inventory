package repository

import (
	"fmt"
	"github.com/alifirfandi/simple-cashier-inventory/entity"
	"github.com/alifirfandi/simple-cashier-inventory/model"
	"gorm.io/gorm"
	"os"
	"strconv"
	"time"
)

type UserRepositoryImpl struct {
	Mysql gorm.DB
}

func mapUserRequestToUserEntity(userEntity *entity.User, userRequest model.UserRequest) {
	userEntity.Name = userRequest.Name
	userEntity.Email = userRequest.Email
	userEntity.Password = userRequest.Password
	userEntity.Role = userRequest.Role
}

func mapUserEntityToUserResponse(userResponse *model.UserResponse, userEntity entity.User) {
	userResponse.Id = userEntity.Id
	userResponse.Name = userEntity.Name
	userResponse.Email = userEntity.Email
	userResponse.Password = userEntity.Password
	userResponse.Role = userEntity.Role
	userResponse.CreatedAt = userEntity.CreatedAt
	userResponse.UpdatedAt = userEntity.UpdateAt
}

func (Repository UserRepositoryImpl) InsertUser(Request model.UserRequest) (Response model.UserResponse, Error error) {
	var user entity.User
	mapUserRequestToUserEntity(&user, Request)
	if Error = Repository.Mysql.Create(&user).Error; Error != nil {
		return Response, Error
	}
	mapUserEntityToUserResponse(&Response, user)
	return Response, Error
}

func (Repository UserRepositoryImpl) GetAllUser(QueryParams model.UserSelectQuery) (Response []model.UserResponse, TotalData int64, Error error) {
	var users []entity.User
	var TotalDataCount int64

	QueryParams.Q = fmt.Sprintf("%%%s%%", QueryParams.Q)

	limirPerPage, _ := strconv.Atoi(os.Getenv("LIMIT_PAGE"))
	QueryParams.Page = (QueryParams.Page - 1) * limirPerPage

	if Error = Repository.Mysql.Table("users").Where("deleted_at IS NULL AND (name LIKE ? OR email LIKE ?)", QueryParams.Q, QueryParams.Q).Count(&TotalDataCount).Error; Error != nil {
		return Response, 0, Error
	}

	if TotalDataCount != 0 {
		if Error = Repository.Mysql.Limit(10).Offset(QueryParams.Page).Where("deleted_at IS NULL AND (name LIKE ? OR email LIKE ?)", QueryParams.Q, QueryParams.Q).Find(&users).Error; Error != nil {
			return Response, 0, Error
		}
	}

	Response = make([]model.UserResponse, len(users))
	for i, product := range users {
		mapUserEntityToUserResponse(&Response[i], product)
	}
	return Response, TotalDataCount, Error
}

func (Repository UserRepositoryImpl) GetUserById(Id int64) (Response model.UserResponse, Error error) {
	var user entity.User
	if Error = Repository.Mysql.Where("deleted_at IS NULL").First(&user, Id).Error; Error != nil {
		return Response, Error
	}
	mapUserEntityToUserResponse(&Response, user)
	return Response, Error
}

func (Repository UserRepositoryImpl) UpdateUser(Id int64, Request model.UserRequest) (Response model.UserResponse, Error error) {
	var newUser entity.User
	mapUserRequestToUserEntity(&newUser, Request)

	var user entity.User
	if Error = Repository.Mysql.Where("deleted_at IS NULL").First(&user, Id).Error; Error != nil {
		return Response, Error
	}

	user.Name = newUser.Name
	user.Email = newUser.Email
	user.Role = newUser.Role
	if newUser.Password != "" {
		user.Password = newUser.Password
	}

	if Error = Repository.Mysql.Save(&user).Error; Error != nil {
		return Response, Error
	}

	mapUserEntityToUserResponse(&Response, user)
	return Response, Error
}

func (Repository UserRepositoryImpl) DeleteUser(Id int64) (Error error) {
	var user entity.User
	if Error = Repository.Mysql.Where("deleted_at IS NULL").First(&user, Id).Error; Error != nil {
		return Error
	}

	if Error = Repository.Mysql.Model(&user).Update("deleted_at", time.Now()).Error; Error != nil {
		return Error
	}
	return Error
}

func NewUserRepository(Mysql *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		Mysql: *Mysql,
	}
}
