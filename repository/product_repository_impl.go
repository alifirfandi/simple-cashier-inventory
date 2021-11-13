package repository

import (
	"github.com/alifirfandi/simple-cashier-inventory/entity"
	"github.com/alifirfandi/simple-cashier-inventory/model"
	"gorm.io/gorm"
)

type ProductRepositoryImpl struct {
	Mysql gorm.DB
}

func NewProductRepository(Mysql *gorm.DB) ProductRepository {
	return &ProductRepositoryImpl{
		Mysql: *Mysql,
	}
}

func (Repository ProductRepositoryImpl) InsertProduct(Request model.ProductRequest) (Response model.ProductResponse, Error error) {
	product := entity.Product{
		Name:     Request.Name,
		Price:    Request.Price,
		Stock:    Request.Stock,
		ImageUrl: Request.ImageUrl,
	}
	Error = Repository.Mysql.Create(&product).Error
	if Error != nil {
		return Response, Error
	}

	Response.Id = product.Id
	Response.Name = product.Name
	Response.Stock = product.Stock
	Response.Price = product.Price
	Response.ImageUrl = product.ImageUrl
	return
}
