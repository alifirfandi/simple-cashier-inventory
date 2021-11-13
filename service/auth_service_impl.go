package service

import (
	"time"

	"github.com/alifirfandi/simple-cashier-inventory/helper"
	"github.com/alifirfandi/simple-cashier-inventory/model"
	"github.com/alifirfandi/simple-cashier-inventory/repository"
	"github.com/alifirfandi/simple-cashier-inventory/validation"

	"github.com/dgrijalva/jwt-go"
)

type AuthServiceImpl struct {
	AuthRepository repository.AuthRepository
}

func NewAuthService(AuthRepository *repository.AuthRepository) AuthService {
	return &AuthServiceImpl{
		AuthRepository: *AuthRepository,
	}
}

func (Service AuthServiceImpl) Login(Request model.AuthRequest) (Response model.AuthResponse, Verified bool, Error error) {
	if Error = validation.LoginValidation(Request); Error != nil {
		return Response, Verified, Error
	}
	Response, userExists, Error := Service.AuthRepository.Login(Request)
	if userExists {
		if helper.CompareHash(Response.Password, Request.Password) {
			accessToken := helper.SignJWT(jwt.MapClaims{
				"exp":   time.Now().Add(24 * time.Hour).Unix(),
				"id":    Response.Id,
				"email": Response.Email,
				"role":  Response.Role,
			})
			Response.AccessToken = accessToken
			return Response, true, Error
		}
	}
	return Response, Verified, Error
}
