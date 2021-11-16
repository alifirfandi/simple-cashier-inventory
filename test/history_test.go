package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alifirfandi/simple-cashier-inventory/config"
	"github.com/alifirfandi/simple-cashier-inventory/controller"
	"github.com/alifirfandi/simple-cashier-inventory/entity"
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
		{
			testName: "Get All Histories Success, With Query '?q=ORD123'",
			setup: func(req *http.Request) {
				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/history?q=ORD123",
			expectedStatusCode: http.StatusOK,
		},
		{
			testName: "Get All Histories Success, With Query '?start_date=2021-11-05'",
			setup: func(req *http.Request) {
				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/history?start_date=2021-11-05",
			expectedStatusCode: http.StatusOK,
		},
		{
			testName: "Get All Histories Success, With Query '?end_date=2021-11-05'",
			setup: func(req *http.Request) {
				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/history?end_date=2021-11-05",
			expectedStatusCode: http.StatusOK,
		},
		{
			testName: "Get All Histories Success, With Query '?page=1'",
			setup: func(req *http.Request) {
				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/history?page=1",
			expectedStatusCode: http.StatusOK,
		},
		{
			testName: "Get All Histories Success, With Query '?q=ORD124page=1&start_date=2021-11-05&end_date=2021-11-05'",
			setup: func(req *http.Request) {
				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/history?q=ORD124page=1&start_date=2021-11-05&end_date=2021-11-05",
			expectedStatusCode: http.StatusOK,
		},
		{
			testName: "Get All Histories Fail, Token Role Not Allowed",
			setup: func(req *http.Request) {
				res := register(gormDb, registerRequest{Email: "email1@email.com", Name: "Name", Password: "password"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/history",
			expectedStatusCode: http.StatusForbidden,
		},
		{
			testName:           "Get All Histories Fail, No Token Provided",
			setup:              func(req *http.Request) {},
			method:             "GET",
			url:                "/history",
			expectedStatusCode: http.StatusUnauthorized,
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
				products := entity.Transaction{
					Invoice:    "ORD0",
					GrandTotal: 1,
					AdminId:    1,
				}
				gormDb.Create(&products)

				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/history/ORD0",
			expectedStatusCode: http.StatusOK,
		},
		{
			testName: "Get History Detail Fail, History Detail Not Found",
			setup: func(req *http.Request) {
				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/history/ORD999",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			testName: "Get History Detail Fail, User Not Allowed",
			setup: func(req *http.Request) {
				res := register(gormDb, registerRequest{Email: "email1@email.com", Name: "Name", Password: "password"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/history/ORD999",
			expectedStatusCode: http.StatusForbidden,
		},
		{
			testName:           "Get History Detail Fail, No Token Provided",
			setup:              func(req *http.Request) {},
			method:             "GET",
			url:                "/history/ORD999",
			expectedStatusCode: http.StatusUnauthorized,
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

func TestExportHistoryPdfDetail(t *testing.T) {
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
			testName: "Export History Pdf Success",
			setup: func(req *http.Request) {
				products := entity.Transaction{
					Invoice:    "ORD0",
					GrandTotal: 1,
					AdminId:    1,
				}
				gormDb.Create(&products)

				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/history/ORD0",
			expectedStatusCode: http.StatusOK,
		},
		{
			testName: "Export History Pdf Fail, History Detail Not Found",
			setup: func(req *http.Request) {
				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/history/ORD999",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			testName: "Export History Pdf Fail, User Not Allowed",

			setup: func(req *http.Request) {
				res := register(gormDb, registerRequest{Email: "email1@email.com", Name: "Name", Password: "password"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/history/ORD999",
			expectedStatusCode: http.StatusForbidden,
		},
		{
			testName:           "Export History Pdf Fail, No Token Provided",
			setup:              func(req *http.Request) {},
			method:             "GET",
			url:                "/history/ORD999",
			expectedStatusCode: http.StatusUnauthorized,
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
