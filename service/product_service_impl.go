package service

import (
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"github.com/alifirfandi/simple-cashier-inventory/helper"
	"github.com/alifirfandi/simple-cashier-inventory/validation"

	"github.com/alifirfandi/simple-cashier-inventory/model"
	"github.com/alifirfandi/simple-cashier-inventory/repository"
)

type ProductServiceImpl struct {
	ProductRepository repository.ProductRepository
}

func (Service ProductServiceImpl) InsertProduct(Request model.ProductRequest, File *multipart.FileHeader) (Response model.ProductResponse, Error error) {
	if Error = validation.ProductValidation(Request); Error != nil {
		return Response, Error
	}

	if File != nil {
		Request.ImageUrl = fmt.Sprintf("%s-%d-%s", helper.RandString(9), time.Now().Unix(), File.Filename)
		if Error = helper.SaveFile(File, os.Getenv("FILE_LOCATION"), Request.ImageUrl); Error != nil {
			return Response, Error
		}
	}

	Response, Error = Service.ProductRepository.InsertProduct(Request)
	return Response, Error
}

func (Service ProductServiceImpl) GetAllProducts() (Response model.ProductListResponse, Error error) {
	products, Error := Service.ProductRepository.GetAllProducts()
	Response = model.ProductListResponse{
		Products: products,
	}

	return Response, Error
}

func (Service ProductServiceImpl) GetProductById(Id int64) (Response model.ProductResponse, Error error) {
	Response, Error = Service.ProductRepository.GetProductById(Id)
	return Response, Error
}

func (Service ProductServiceImpl) UpdateProductById(Id int64, Request model.ProductRequest, File *multipart.FileHeader) (Response model.ProductResponse, Error error) {
	if Error = validation.ProductValidation(Request); Error != nil {
		return Response, Error
	}

	if File != nil {
		Request.ImageUrl = fmt.Sprintf("%s-%d-%s", helper.RandString(9), time.Now().Unix(), File.Filename)
		if Error = helper.SaveFile(File, os.Getenv("FILE_LOCATION"), Request.ImageUrl); Error != nil {
			return Response, Error
		}
	}

	Response, Error = Service.ProductRepository.UpdateProductById(Id, Request)
	return Response, Error
}

func (Service ProductServiceImpl) DeleteProductById(Id int64) (Error error) {

	if Error = Service.ProductRepository.DeleteProductById(Id); Error != nil {
		return Error
	}

	return nil
}

func NewProductService(ProductRepository *repository.ProductRepository) ProductService {
	return &ProductServiceImpl{
		ProductRepository: *ProductRepository,
	}
}
