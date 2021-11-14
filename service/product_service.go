package service

import (
	"mime/multipart"

	"github.com/alifirfandi/simple-cashier-inventory/model"
)

type ProductService interface {
	InsertProduct(Request model.ProductRequest, File *multipart.FileHeader) (Response model.ProductResponse, Error error)
	GetAllProducts() (Response model.ProductListResponse, Error error)
	GetProductById(Id int64) (Response model.ProductResponse, Error error)
	UpdateProductById(Id int64, Request model.ProductRequest, File *multipart.FileHeader) (Response model.ProductResponse, Error error)
	DeleteProductById(Id int64) (Error error)
}
