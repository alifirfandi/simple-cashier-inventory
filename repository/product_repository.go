package repository

import "github.com/alifirfandi/simple-cashier-inventory/model"

type ProductRepository interface {
	InsertProduct(Request model.ProductRequest) (Response model.ProductResponse, Error error)
	GetAllProducts(Query model.ProductSelectQuery) (Response []model.ProductResponse, Error error)
	GetProductById(Id int64) (Response model.ProductResponse, Error error)
	UpdateProductById(Id int64, Request model.ProductRequest) (Response model.ProductResponse, Error error)
	DeleteProductById(Id int64) (Error error)
	CountProducts() (Result int64, Error error)
}
