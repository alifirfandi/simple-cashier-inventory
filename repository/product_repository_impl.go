package repository

import (
	"fmt"
	"time"

	"github.com/alifirfandi/simple-cashier-inventory/entity"
	"github.com/alifirfandi/simple-cashier-inventory/model"
	"gorm.io/gorm"
)

func mapProductRequestToProductEntity(productEntity *entity.Product, productRequest model.ProductRequest) {
	productEntity.Name = productRequest.Name
	productEntity.ImageUrl = productRequest.ImageUrl
	productEntity.Price = productRequest.Price
	productEntity.Stock = productRequest.Stock
}

func mapProductEntityToProductResponse(productResponse *model.ProductResponse, productEntity entity.Product) {
	productResponse.Id = productEntity.Id
	productResponse.Name = productEntity.Name
	productResponse.Stock = productEntity.Stock
	productResponse.Price = productEntity.Price
	productResponse.ImageUrl = productEntity.ImageUrl
	productResponse.CreatedAt = productEntity.CreatedAt.Format(time.RFC3339)
	productResponse.UpdatedAt = productEntity.UpdateAt.Format(time.RFC3339)
}

type ProductRepositoryImpl struct {
	Mysql gorm.DB
}

func (Repository ProductRepositoryImpl) InsertProduct(Request model.ProductRequest) (Response model.ProductResponse, Error error) {
	var product entity.Product
	mapProductRequestToProductEntity(&product, Request)
	if Error = Repository.Mysql.Create(&product).Error; Error != nil {
		return Response, Error
	}
	mapProductEntityToProductResponse(&Response, product)
	return Response, Error
}

func (Repository ProductRepositoryImpl) GetAllProducts(Query model.ProductSelectQuery) (Response []model.ProductResponse, Error error) {
	var products []entity.Product
	Error = Repository.Mysql.Where("name LIKE ? AND deleted_at IS NULL", Query.Search).Order(fmt.Sprintf("%s %s", Query.SortField, Query.SortBy)).Limit(Query.Limit).Offset(Query.Start).Find(&products).Error
	if Error != nil {
		return Response, Error
	}

	Response = make([]model.ProductResponse, len(products))
	for i, product := range products {
		mapProductEntityToProductResponse(&Response[i], product)
	}
	return Response, Error
}

func (Repository ProductRepositoryImpl) GetProductById(Id int64) (Response model.ProductResponse, Error error) {
	var product entity.Product
	if Error = Repository.Mysql.Where("deleted_at IS NULL").First(&product, Id).Error; Error != nil {
		return Response, Error
	}
	mapProductEntityToProductResponse(&Response, product)
	return Response, Error
}

func (Repository ProductRepositoryImpl) UpdateProductById(Id int64, Request model.ProductRequest) (Response model.ProductResponse, Error error) {
	var newProduct entity.Product
	mapProductRequestToProductEntity(&newProduct, Request)

	var product entity.Product
	if Error = Repository.Mysql.Where("deleted_at IS NULL").First(&product, Id).Error; Error != nil {
		return Response, Error
	}

	product.Name = newProduct.Name
	product.Stock = newProduct.Stock
	product.Price = newProduct.Price
	if newProduct.ImageUrl != "" {
		product.ImageUrl = newProduct.ImageUrl
	}

	if Error = Repository.Mysql.Save(&product).Error; Error != nil {
		return Response, Error
	}

	mapProductEntityToProductResponse(&Response, product)
	return Response, Error
}

func (Repository ProductRepositoryImpl) DeleteProductById(Id int64) (Error error) {
	var product entity.Product
	if Error = Repository.Mysql.Where("deleted_at IS NULL").First(&product, Id).Error; Error != nil {
		return Error
	}

	if Error = Repository.Mysql.Model(&product).Update("deleted_at", time.Now()).Error; Error != nil {
		return Error
	}
	return Error
}

func (Repository ProductRepositoryImpl) CountProducts() (Result int64, Error error) {
	Error = Repository.Mysql.Model(&entity.Product{}).Count(&Result).Error
	return Result, Error
}

func NewProductRepository(Mysql *gorm.DB) ProductRepository {
	return &ProductRepositoryImpl{
		Mysql: *Mysql,
	}
}
