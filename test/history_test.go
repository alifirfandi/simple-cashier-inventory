package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alifirfandi/simple-cashier-inventory/config"
	"github.com/alifirfandi/simple-cashier-inventory/controller"
	"github.com/alifirfandi/simple-cashier-inventory/model"
	"github.com/alifirfandi/simple-cashier-inventory/repository"
	"github.com/alifirfandi/simple-cashier-inventory/service"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func historyInit(gormDb *gorm.DB) (app *fiber.App, historyService service.HistoryService, historyRepository repository.HistoryRepository, authService service.AuthService) {
	authRepository := repository.NewAuthRepository(gormDb)
	authService = service.NewAuthService(&authRepository)
	historyRepository = repository.NewHistoryRepository(gormDb)
	historyService = service.NewHistoryService(&historyRepository)
	historyController := controller.NewHistoryController(&historyService)

	app = fiber.New(config.NewFiberConfig())
	historyController.Route(app)
	return
}

func TestGetAllHistories(t *testing.T) {
	gormDb := dbInit()
	app, _, _, authService := historyInit(gormDb)
	clearDb(gormDb)

	testCases := []struct {
		testName           string
		setup              func(req *http.Request)
		method             string
		url                string
		expectedStatusCode int
	}{
		{
			testName: "Get All Histories Success",
			setup: func(req *http.Request) {
				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/history",
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			req := httptest.NewRequest(testCase.method, testCase.url, nil)
			testCase.setup(req)

			resp, _ := app.Test(req)
			if resp.StatusCode != testCase.expectedStatusCode {
				body, _ := ioutil.ReadAll(resp.Body)
				t.Fatal(string(body))
			}
		})
	}
}

func TestGetHistoryDetail(t *testing.T) {
	gormDb := dbInit()
	app, _, _, authService := historyInit(gormDb)
	clearDb(gormDb)

	testCases := []struct {
		testName           string
		setup              func(req *http.Request)
		method             string
		url                string
		expectedStatusCode int
	}{
		{
			testName: "Get History Detail Success",
			setup: func(req *http.Request) {
				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/history/1",
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			req := httptest.NewRequest(testCase.method, testCase.url, nil)
			testCase.setup(req)

			resp, _ := app.Test(req)
			if resp.StatusCode != testCase.expectedStatusCode {
				body, _ := ioutil.ReadAll(resp.Body)
				t.Fatal(string(body))
			}
		})
	}
}
