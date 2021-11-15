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

	"github.com/alifirfandi/simple-cashier-inventory/entity"

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
			testName: "Insert Product Success, With Image",
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
				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
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
				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
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
				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
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
				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
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
				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
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
				res := register(gormDb, registerRequest{Email: "email6@email.com", Name: "Name", Password: "password"})
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

func TestGetAllProducts(t *testing.T) {
	gormDb := dbInit()
	app, _, _, authService := productInit(gormDb)
	clearDb(gormDb)

	testCases := []struct {
		testName           string
		setup              func(req *http.Request)
		method             string
		url                string
		expectedStatusCode int
	}{
		{
			testName: "Get All Product Success",
			setup: func(req *http.Request) {
				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/product",
			expectedStatusCode: http.StatusOK,
		},
		{
			testName: "Get All Product Success",
			setup: func(req *http.Request) {
				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/product",
			expectedStatusCode: http.StatusOK,
		},
		{
			testName: "Get All Product Success, With Query '?q=Product%20Name'",
			setup: func(req *http.Request) {
				products := make([]entity.Product, 15)
				for i := 0; i < len(products); i++ {
					products[i] = entity.Product{
						Name:  "Product Name",
						Price: 1,
						Stock: 1,
					}
				}
				gormDb.Create(&products)

				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/product?q=Product%20Name",
			expectedStatusCode: http.StatusOK,
		},
		{
			testName: "Get All Product Success, With Query '?page=1'",
			setup: func(req *http.Request) {
				products := make([]entity.Product, 15)
				for i := 0; i < len(products); i++ {
					products[i] = entity.Product{
						Name:  "Product Name",
						Price: 1,
						Stock: 1,
					}
				}
				gormDb.Create(&products)

				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/product?page=1",
			expectedStatusCode: http.StatusOK,
		},
		{
			testName: "Get All Product Success, With Query '?sort=name_asc'",
			setup: func(req *http.Request) {
				products := make([]entity.Product, 15)
				for i := 0; i < len(products); i++ {
					products[i] = entity.Product{
						Name:  "Product Name",
						Price: 1,
						Stock: 1,
					}
				}
				gormDb.Create(&products)

				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/product?sort=name_asc",
			expectedStatusCode: http.StatusOK,
		},
		{
			testName: "Get All Product Success, With Query '?sort=name_desc'",
			setup: func(req *http.Request) {
				products := make([]entity.Product, 15)
				for i := 0; i < len(products); i++ {
					products[i] = entity.Product{
						Name:  "Product Name",
						Price: 1,
						Stock: 1,
					}
				}
				gormDb.Create(&products)

				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/product?sort=name_desc",
			expectedStatusCode: http.StatusOK,
		},
		{
			testName: "Get All Product Success, With Query '?sort=created_at_asc'",
			setup: func(req *http.Request) {
				products := make([]entity.Product, 15)
				for i := 0; i < len(products); i++ {
					products[i] = entity.Product{
						Name:  "Product Name",
						Price: 1,
						Stock: 1,
					}
				}
				gormDb.Create(&products)

				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/product?sort=created_at_asc",
			expectedStatusCode: http.StatusOK,
		},
		{
			testName: "Get All Product Success, With Query '?sort=created_at_desc'",
			setup: func(req *http.Request) {
				products := make([]entity.Product, 15)
				for i := 0; i < len(products); i++ {
					products[i] = entity.Product{
						Name:  "Product Name",
						Price: 1,
						Stock: 1,
					}
				}
				gormDb.Create(&products)

				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/product?sort=created_at_desc",
			expectedStatusCode: http.StatusOK,
		},
		{
			testName: "Get All Product Success, With Query '?sort=updated_at_asc'",
			setup: func(req *http.Request) {
				products := make([]entity.Product, 15)
				for i := 0; i < len(products); i++ {
					products[i] = entity.Product{
						Name:  "Product Name",
						Price: 1,
						Stock: 1,
					}
				}
				gormDb.Create(&products)

				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/product?sort=updated_at_desc",
			expectedStatusCode: http.StatusOK,
		},
		{
			testName: "Get All Product Success, With Query '?sort=price_asc'",
			setup: func(req *http.Request) {
				products := make([]entity.Product, 15)
				for i := 0; i < len(products); i++ {
					products[i] = entity.Product{
						Name:  "Product Name",
						Price: 1,
						Stock: 1,
					}
				}
				gormDb.Create(&products)

				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/product?sort=price_asc",
			expectedStatusCode: http.StatusOK,
		},
		{
			testName: "Get All Product Success, With Query '?sort=price_desc'",
			setup: func(req *http.Request) {
				products := make([]entity.Product, 15)
				for i := 0; i < len(products); i++ {
					products[i] = entity.Product{
						Name:  "Product Name",
						Price: 1,
						Stock: 1,
					}
				}
				gormDb.Create(&products)

				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/product?sort=price_desc",
			expectedStatusCode: http.StatusOK,
		},
		{
			testName:           "Get All Product Fail, No Token Provided",
			setup:              func(req *http.Request) {},
			method:             "GET",
			url:                "/product",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			testName: "Get All Product Fail, Token Is Not SUPERADMIN",
			setup: func(req *http.Request) {
				res := register(gormDb, registerRequest{Email: "email1@email.com", Name: "Name", Password: "password"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/product",
			expectedStatusCode: http.StatusForbidden,
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

func TestGetProductDetail(t *testing.T) {
	gormDb := dbInit()
	app, _, _, authService := productInit(gormDb)
	clearDb(gormDb)

	testCases := []struct {
		testName           string
		setup              func(req *http.Request)
		method             string
		url                string
		expectedStatusCode int
	}{
		{
			testName: "Get Product Detail Success",
			setup: func(req *http.Request) {
				product := entity.Product{
					Id:    1,
					Name:  "Product Name",
					Price: 1,
					Stock: 1,
				}
				gormDb.Create(&product)

				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/product/1",
			expectedStatusCode: http.StatusOK,
		},
		{
			testName: "Get Product Detail Fail, Id Not Found",
			setup: func(req *http.Request) {
				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/product/999",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			testName:           "Get Product Detail Fail, No Token Provided",
			setup:              func(req *http.Request) {},
			method:             "GET",
			url:                "/product/999",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			testName: "Get Product Detail Fail, Token Is Not SUPERADMIN",
			setup: func(req *http.Request) {
				res := register(gormDb, registerRequest{Email: "email1@email.com", Name: "Name", Password: "password"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "GET",
			url:                "/product/999",
			expectedStatusCode: http.StatusForbidden,
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

func TestUpdateProduct(t *testing.T) {
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
			testName: "Update Product Success, With Image",
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
				product := entity.Product{
					Id:    1,
					Name:  "Product Name",
					Price: 1,
					Stock: 1,
				}
				gormDb.Create(&product)

				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method: "PUT",
			url:    "/product/1",
			fields: map[string]string{
				"name":  "New Product Name",
				"stock": "10",
				"price": "1000",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			testName: "Update Product Fail, Name Not Provided",
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
				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method: "PUT",
			url:    "/product/999",
			fields: map[string]string{
				"stock": "10",
				"price": "1000",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			testName: "Update Product Fail, Stock Not Provided",
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
				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method: "PUT",
			url:    "/product/999",
			fields: map[string]string{
				"name":  "New Product Name",
				"price": "1000",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			testName: "Update Product Fail, Price Not Provided",
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
				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method: "PUT",
			url:    "/product/999",
			fields: map[string]string{
				"name":  "New Product Name",
				"stock": "10",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			testName: "Update Product Fail, No Data Provided",
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
				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "PUT",
			url:                "/product/999",
			fields:             map[string]string{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			testName: "Update Product Fail, No Token Provided",
			initForm: func(fields map[string]string) (*bytes.Buffer, *multipart.Writer) {
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)

				err := writer.Close()
				if err != nil {
					panic(err)
				}

				return body, writer
			},
			setup:              func(req *http.Request) {},
			method:             "PUT",
			url:                "/product/999",
			fields:             map[string]string{},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			testName: "Update Product Fail, Role Not Allowed",
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
				res := register(gormDb, registerRequest{Email: "email6@email.com", Name: "Name", Password: "password"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "PUT",
			url:                "/product/999",
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

func TestDeleteProduct(t *testing.T) {
	gormDb := dbInit()
	app, _, _, authService := productInit(gormDb)
	clearDb(gormDb)

	testCases := []struct {
		testName           string
		setup              func(req *http.Request)
		method             string
		url                string
		expectedStatusCode int
	}{
		{
			testName: "Delete Product Success",
			setup: func(req *http.Request) {
				product := entity.Product{
					Id:    1,
					Name:  "Product Name",
					Price: 1,
					Stock: 1,
				}
				gormDb.Create(&product)

				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "DELETE",
			url:                "/product/1",
			expectedStatusCode: http.StatusOK,
		},
		{
			testName: "Delete Product Fail, Id Not Found",
			setup: func(req *http.Request) {
				res, _, _ := authService.Login(model.AuthRequest{Email: "superadmin@golang.com", Password: "12345678"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "DELETE",
			url:                "/product/999",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			testName:           "Delete Product Fail, No Token Provided",
			setup:              func(req *http.Request) {},
			method:             "DELETE",
			url:                "/product/999",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			testName: "Delete Product Fail, Token Is Not SUPERADMIN",
			setup: func(req *http.Request) {
				res := register(gormDb, registerRequest{Email: "email8@email.com", Name: "Name", Password: "password"})
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
			},
			method:             "DELETE",
			url:                "/product/999",
			expectedStatusCode: http.StatusForbidden,
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
