package test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alifirfandi/simple-cashier-inventory/entity"

	"github.com/alifirfandi/simple-cashier-inventory/config"
	"github.com/alifirfandi/simple-cashier-inventory/controller"
	"github.com/alifirfandi/simple-cashier-inventory/helper"
	"github.com/alifirfandi/simple-cashier-inventory/repository"
	"github.com/alifirfandi/simple-cashier-inventory/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type registerRequest struct {
	Name     string
	Email    string
	Password string
}

type registerResponse struct {
	Id          int64
	Name        string
	Email       string
	Role        string
	AccessToken string
}

func register(gormDb *gorm.DB, Request registerRequest) (Response registerResponse) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(Request.Password), bcrypt.DefaultCost)
	user := entity.User{
		Email:     Request.Email,
		Name:      Request.Name,
		Password:  string(hashed),
		CreatedAt: time.Now(),
		UpdateAt:  time.Now(),
	}
	gormDb.Create(&user)

	Response.Id = user.Id
	Response.Name = user.Name
	Response.Email = user.Email
	Response.Role = user.Role

	accessToken := helper.SignJWT(jwt.MapClaims{
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"id":    Response.Id,
		"email": Response.Email,
		"role":  Response.Role,
	})
	Response.AccessToken = accessToken
	return
}

func authInit(gormDb *gorm.DB) (app *fiber.App, authService service.AuthService, authRepository repository.AuthRepository) {
	authRepository = repository.NewAuthRepository(gormDb)
	authService = service.NewAuthService(&authRepository)
	authController := controller.NewAuthController(&authService)

	app = fiber.New(config.NewFiberConfig())
	authController.Route(app)
	return
}

func TestLogin(t *testing.T) {
	gormDb := dbInit()
	app, _, _ := authInit(gormDb)
	clearDb(gormDb)

	testCases := []struct {
		testName           string
		init               func(req *http.Request)
		method             string
		url                string
		body               string
		expectedStatusCode int
	}{
		{
			testName: "Login Success",
			init: func(req *http.Request) {
				req.Header.Add("Content-Type", "application/json")
				user := registerRequest{
					Email:    "email5@email.com",
					Name:     "Name",
					Password: "password",
				}
				register(gormDb, user)
			},
			method:             "POST",
			url:                "/auth/login",
			body:               `{"email": "email5@email.com", "password": "password"}`,
			expectedStatusCode: http.StatusOK,
		},
		{
			testName: "Login Fail, User Not Registered",
			init: func(req *http.Request) {
				req.Header.Add("Content-Type", "application/json")
			},
			method:             "POST",
			url:                "/auth/login",
			body:               `{"email": "email6@email.com", "password": "password"}`,
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			testName: "Login Fail, Wrong Password",
			init: func(req *http.Request) {
				req.Header.Add("Content-Type", "application/json")
				user := registerRequest{
					Email:    "email7@email.com",
					Name:     "Name",
					Password: "password",
				}
				register(gormDb, user)
			},
			method:             "POST",
			url:                "/auth/login",
			body:               `{"email": "email7@email.com", "password": "wrong_password"}`,
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			req := httptest.NewRequest(testCase.method, testCase.url, strings.NewReader(testCase.body))
			testCase.init(req)

			resp, _ := app.Test(req)
			if resp.StatusCode != testCase.expectedStatusCode {
				body, _ := ioutil.ReadAll(resp.Body)
				t.Fatal(string(body))
			}
		})
	}
}
