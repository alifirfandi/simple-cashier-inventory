package repository

import "github.com/alifirfandi/simple-cashier-inventory/model"

type ProductRepository interface {
	InsertProduct(Request model.ProductRequest) (Response model.ProductResponse, Error error)
}
