package service

import (
	"fmt"
	"math"
	"mime/multipart"
	"os"
	"strconv"
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

func (Service ProductServiceImpl) GetAllProducts(Query model.ProductRequestQuery) (Response model.ProductListResponse, Error error) {
	limit, Error := strconv.Atoi(os.Getenv("LIMIT_PER_PAGE"))
	if Error != nil {
		return Response, Error
	}

	sort, sortBy := helper.SplitLastStr(Query.Sort, "_")
	page := (Query.Page - 1) * limit

	products, Error := Service.ProductRepository.GetAllProducts(model.ProductSelectQuery{
		Search:    fmt.Sprintf(`%%%s%%`, Query.Q),
		SortField: sort,
		SortBy:    sortBy,
		Start:     page,
		Limit:     limit,
	})
	if Error != nil {
		return Response, Error
	}

	productsCount, Error := Service.ProductRepository.CountProducts()
	if Error != nil {
		return Response, Error
	}

	Response = model.ProductListResponse{
		TotalData:    int(productsCount),
		TotalPage:    int(math.Ceil(float64(productsCount) / float64(limit))),
		CurrentPage:  Query.Page,
		LimitPerPage: limit,
		Products:     products,
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
