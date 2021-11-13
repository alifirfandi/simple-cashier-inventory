package service

import (
	"mime/multipart"
	"os"

	"github.com/alifirfandi/simple-cashier-inventory/helper"
	"github.com/alifirfandi/simple-cashier-inventory/validation"

	"github.com/alifirfandi/simple-cashier-inventory/model"
	"github.com/alifirfandi/simple-cashier-inventory/repository"
)

type ProductServiceImpl struct {
	ProductRepository repository.ProductRepository
}

func NewProductService(ProductRepository *repository.ProductRepository) ProductService {
	return &ProductServiceImpl{
		ProductRepository: *ProductRepository,
	}
}

func (Service ProductServiceImpl) InsertProduct(Request model.ProductRequest, File *multipart.FileHeader) (Response model.ProductResponse, Error error) {
	if Error = validation.ProductValidation(Request); Error != nil {
		return Response, Error
	}

	if File != nil {
		Request.ImageUrl = File.Filename

		if Error = helper.SaveFile(os.Getenv("FILE_LOCATION"), File); Error != nil {
			return Response, Error
		}
	}

	Response, Error = Service.ProductRepository.InsertProduct(Request)
	return Response, Error
}
