package service

import (
	"mime/multipart"

	"github.com/alifirfandi/simple-cashier-inventory/model"
)

type ProductService interface {
	InsertProduct(Request model.ProductRequest, File *multipart.FileHeader) (Response model.ProductResponse, Error error)
}
