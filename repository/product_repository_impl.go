package repository

import (
	"time"

	"github.com/alifirfandi/simple-cashier-inventory/entity"
	"github.com/alifirfandi/simple-cashier-inventory/model"
	"gorm.io/gorm"
)

func mapProductRequestToEntity(productEntity *entity.Product, productRequest model.ProductRequest) {
	productEntity.Name = productRequest.Name
	productEntity.ImageUrl = productRequest.ImageUrl
	productEntity.Price = productRequest.Price
	productEntity.Stock = productRequest.Stock
}

func mapProductEntityToResponse(productResponse *model.ProductResponse, productEntity entity.Product) {
	productResponse.Id = productEntity.Id
	productResponse.Name = productEntity.Name
	productResponse.Stock = productEntity.Stock
	productResponse.Price = productEntity.Price
	productResponse.ImageUrl = productEntity.ImageUrl
	productResponse.CreatedAt = productEntity.CreatedAt
	productResponse.UpdatedAt = productEntity.UpdateAt
}

type ProductRepositoryImpl struct {
	Mysql gorm.DB
}

func (Repository ProductRepositoryImpl) InsertProduct(Request model.ProductRequest) (Response model.ProductResponse, Error error) {
	var product entity.Product
	mapProductRequestToEntity(&product, Request)
	if Error = Repository.Mysql.Create(&product).Error; Error != nil {
		return Response, Error
	}
	mapProductEntityToResponse(&Response, product)
	return Response, Error
}

func (Repository ProductRepositoryImpl) GetAllProducts() (Response []model.ProductResponse, Error error) {
	var products []entity.Product
	if Error = Repository.Mysql.Where("deleted_at IS NULL").Find(&products).Error; Error != nil {
		return Response, Error
	}

	Response = make([]model.ProductResponse, len(products))
	for i, product := range products {
		mapProductEntityToResponse(&Response[i], product)
		//Response[i] = model.ProductResponse{
		//	Id:        product.Id,
		//	Name:      product.Name,
		//	ImageUrl:  product.ImageUrl,
		//	Price:     product.Price,
		//	Stock:     product.Stock,
		//	CreatedAt: product.CreatedAt,
		//	UpdatedAt: product.UpdateAt,
		//}
	}
	return Response, Error
}

func (Repository ProductRepositoryImpl) GetProductById(Id int64) (Response model.ProductResponse, Error error) {
	var product entity.Product
	if Error = Repository.Mysql.Where("deleted_at IS NULL").First(&product, Id).Error; Error != nil {
		return Response, Error
	}
	mapProductEntityToResponse(&Response, product)
	return Response, Error
}

func (Repository ProductRepositoryImpl) UpdateProductById(Id int64, Request model.ProductRequest) (Response model.ProductResponse, Error error) {
	var newProduct entity.Product
	mapProductRequestToEntity(&newProduct, Request)

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

	mapProductEntityToResponse(&Response, product)
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

func NewProductRepository(Mysql *gorm.DB) ProductRepository {
	return &ProductRepositoryImpl{
		Mysql: *Mysql,
	}
}
