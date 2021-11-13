package test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alifirfandi/simple-cashier-inventory/model"

	"github.com/alifirfandi/simple-cashier-inventory/config"
	"github.com/alifirfandi/simple-cashier-inventory/controller"
	"github.com/alifirfandi/simple-cashier-inventory/repository"
	"github.com/alifirfandi/simple-cashier-inventory/service"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func productInit(gormDb *gorm.DB) (app *fiber.App, productService service.ProductService, productRepository repository.ProductRepository, authService service.AuthService) {
	authRepository := repository.NewAuthRepository(gormDb)
	authService = service.NewAuthService(&authRepository)
	productRepository = repository.NewProductRepository(gormDb)
	productService = service.NewProductService(&productRepository)
	productController := controller.NewProductController(&productService)

	app = fiber.New(config.NewFiberConfig())
	productController.Route(app)
	return
}

// Ref: Testing Go http.Request.FormFile?
// https://stackoverflow.com/a/63115215/12976234
func TestCreateProduct(t *testing.T) {
	gormDb := dbInit()
	app, _, _, authService := productInit(gormDb)
	clearDb(gormDb)

	testCases := []struct {
		testName           string
		initForm           func(map[string]string) (*bytes.Buffer, *multipart.Writer)
		setup              func(req *http.Request)
		fields             map[string]string
		method             string
		url                string
		expectedStatusCode int
	}{
		{
			testName: "Insert Product Success",
			initForm: func(fields map[string]string) (*bytes.Buffer, *multipart.Writer) {
				path := "../fixtures/image.png"

				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)

				for k, v := range fields {
					err := writer.WriteField(k, v)
					if err != nil {
						panic(err)
					}
				}

				part, err := writer.CreateFormFile("image", path)
				if err != nil {
					panic(err)
				}

				sample, err := os.Open(path)
				if err != nil {
					panic(err)
				}

				_, err = io.Copy(part, sample)
				if err != nil {
					panic(err)
				}

				err = writer.Close()
				if err != nil {
					panic(err)
				}

				return body, writer
			},
			setup: func(req *http.Request) {
				register(gormDb, registerRequest{Email: "email1@email.com", Name: "Name", Password: "password"})
				gormDb.Exec("UPDATE users SET role = 'SUPERADMIN' WHERE email = 'email1@email.com'")
				res, _, _ := authService.Login(model.AuthRequest{Email: "email1@email.com", Password: "password"})

				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method: "POST",
			url:    "/product",
			fields: map[string]string{
				"name":  "Product Name",
				"stock": "10",
				"price": "1000",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			testName: "Insert Product Fail, Name Not Provided",
			initForm: func(fields map[string]string) (*bytes.Buffer, *multipart.Writer) {
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)

				for k, v := range fields {
					err := writer.WriteField(k, v)
					if err != nil {
						panic(err)
					}
				}

				err := writer.Close()
				if err != nil {
					panic(err)
				}

				return body, writer
			},
			setup: func(req *http.Request) {
				register(gormDb, registerRequest{Email: "email2@email.com", Name: "Name", Password: "password"})
				gormDb.Exec("UPDATE users SET role = 'SUPERADMIN' WHERE email = 'email2@email.com'")
				res, _, _ := authService.Login(model.AuthRequest{Email: "email2@email.com", Password: "password"})

				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method: "POST",
			url:    "/product",
			fields: map[string]string{
				"stock": "10",
				"price": "1000",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			testName: "Insert Product Fail, Stock Not Provided",
			initForm: func(fields map[string]string) (*bytes.Buffer, *multipart.Writer) {
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)

				for k, v := range fields {
					err := writer.WriteField(k, v)
					if err != nil {
						panic(err)
					}
				}

				err := writer.Close()
				if err != nil {
					panic(err)
				}

				return body, writer
			},
			setup: func(req *http.Request) {
				register(gormDb, registerRequest{Email: "email3@email.com", Name: "Name", Password: "password"})
				gormDb.Exec("UPDATE users SET role = 'SUPERADMIN' WHERE email = 'email3@email.com'")
				res, _, _ := authService.Login(model.AuthRequest{Email: "email3@email.com", Password: "password"})

				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method: "POST",
			url:    "/product",
			fields: map[string]string{
				"name":  "Product Name",
				"price": "1000",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			testName: "Insert Product Fail, Price Not Provided",
			initForm: func(fields map[string]string) (*bytes.Buffer, *multipart.Writer) {
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)

				for k, v := range fields {
					err := writer.WriteField(k, v)
					if err != nil {
						panic(err)
					}
				}

				err := writer.Close()
				if err != nil {
					panic(err)
				}

				return body, writer
			},
			setup: func(req *http.Request) {
				register(gormDb, registerRequest{Email: "email4@email.com", Name: "Name", Password: "password"})
				gormDb.Exec("UPDATE users SET role = 'SUPERADMIN' WHERE email = 'email4@email.com'")
				res, _, _ := authService.Login(model.AuthRequest{Email: "email4@email.com", Password: "password"})

				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method: "POST",
			url:    "/product",
			fields: map[string]string{
				"name":  "Product Name",
				"stock": "10",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			testName: "Insert Product Fail, No Data Provided",
			initForm: func(fields map[string]string) (*bytes.Buffer, *multipart.Writer) {
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)

				err := writer.Close()
				if err != nil {
					panic(err)
				}

				return body, writer
			},
			setup: func(req *http.Request) {
				register(gormDb, registerRequest{Email: "email5@email.com", Name: "Name", Password: "password"})
				gormDb.Exec("UPDATE users SET role = 'SUPERADMIN' WHERE email = 'email5@email.com'")
				res, _, _ := authService.Login(model.AuthRequest{Email: "email5@email.com", Password: "password"})

				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "POST",
			url:                "/product",
			fields:             map[string]string{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			testName: "Insert Product Fail, No Token Provided",
			initForm: func(fields map[string]string) (*bytes.Buffer, *multipart.Writer) {
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)

				err := writer.Close()
				if err != nil {
					panic(err)
				}

				return body, writer
			},
			setup: func(req *http.Request) {
			},
			method:             "POST",
			url:                "/product",
			fields:             map[string]string{},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			testName: "Insert Product Fail, Role Not Allowed",
			initForm: func(fields map[string]string) (*bytes.Buffer, *multipart.Writer) {
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)

				err := writer.Close()
				if err != nil {
					panic(err)
				}

				return body, writer
			},
			setup: func(req *http.Request) {
				register(gormDb, registerRequest{Email: "email6@email.com", Name: "Name", Password: "password"})
				res, _, _ := authService.Login(model.AuthRequest{Email: "email6@email.com", Password: "password"})

				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "POST",
			url:                "/product",
			fields:             map[string]string{},
			expectedStatusCode: http.StatusForbidden,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			body, multi := testCase.initForm(testCase.fields)
			req := httptest.NewRequest(testCase.method, testCase.url, body)
			req.Header.Set("Content-Type", multi.FormDataContentType())
			testCase.setup(req)

			resp, _ := app.Test(req)
			if resp.StatusCode != testCase.expectedStatusCode {
				body, _ := ioutil.ReadAll(resp.Body)
				t.Fatal(string(body))
			}
		})
	}
}
