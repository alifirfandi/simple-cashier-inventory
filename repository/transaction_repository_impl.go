package repository

import (
	"errors"
	"github.com/alifirfandi/simple-cashier-inventory/entity"
	"github.com/alifirfandi/simple-cashier-inventory/helper"
	"github.com/alifirfandi/simple-cashier-inventory/model"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type TransactionRepositoryImpl struct {
	Mysql gorm.DB
}

func mapCartRequestToCartEntity(cartEntity *entity.Cart, cartRequest model.CartRequest) {
	cartEntity.AdminId = cartRequest.AdminId
	cartEntity.ProductId = cartRequest.ProductId
	cartEntity.Qty = cartRequest.Qty
}

func mapCartEntityToCartResponse(cartResponse *model.CartResponse, cartEntity entity.Cart, productEntity entity.Product) {
	cartResponse.Id = cartEntity.Id
	cartResponse.ProductId = productEntity.Id
	cartResponse.Name = productEntity.Name
	cartResponse.ImageUrl = productEntity.ImageUrl
	cartResponse.Price = productEntity.Price
	cartResponse.Qty = cartEntity.Qty
	cartResponse.SubTotal = cartEntity.Qty * productEntity.Price
}

func mapTransactionDetailEntityToTransactionDetailResponse(transactionDetailResponse *model.TransactionDetailResponse, transactionDetailEntity entity.TransactionDetail, productEntity entity.Product) {
	transactionDetailResponse.Id = transactionDetailEntity.Id
	transactionDetailResponse.ProductId = productEntity.Id
	transactionDetailResponse.Name = productEntity.Name
	transactionDetailResponse.ImageUrl = productEntity.ImageUrl
	transactionDetailResponse.Price = productEntity.Price
	transactionDetailResponse.Qty = transactionDetailEntity.Qty
	transactionDetailResponse.SubTotal = transactionDetailEntity.Qty * productEntity.Price
}

func (Repository TransactionRepositoryImpl) InsertCart(Request model.CartRequest) (Response model.CartResponse, Error error) {
	var cart entity.Cart
	var product entity.Product

	mapCartRequestToCartEntity(&cart, Request)
	if Error = Repository.Mysql.Where("deleted_at IS NULL").First(&product, Request.ProductId).Error; Error != nil {
		return Response, Error
	}

	if product.Stock < Request.Qty {
		Error = errors.New("STOCK_UNAVAILABLE")
		return Response, Error
	}

	if Error = Repository.Mysql.Create(&cart).Error; Error != nil {
		return Response, Error
	}
	mapCartEntityToCartResponse(&Response, cart, product)
	return Response, Error
}

func (Repository TransactionRepositoryImpl) GetAllCart(idUser int64) (Response model.CartListResponse, Error error) {
	var listCart []entity.Cart
	var product entity.Product
	var GrandTotal int

	if Error = Repository.Mysql.Where("deleted_at IS NULL AND admin_id = ?", idUser).Find(&listCart).Error; Error != nil {
		return Response, Error
	}

	cartListResponse := make([]model.CartResponse, len(listCart))
	for i, cart := range listCart {
		product.Id = cart.ProductId
		if Error = Repository.Mysql.Where("deleted_at IS NULL").First(&product).Error; Error != nil {
			return Response, Error
		}
		mapCartEntityToCartResponse(&cartListResponse[i], cart, product)
		GrandTotal += cartListResponse[i].SubTotal
	}

	Response.Details = cartListResponse
	Response.GrandTotal = GrandTotal
	return Response, Error
}

func (Repository TransactionRepositoryImpl) UpdateCart(Id int64, Request model.CartRequest) (Response model.CartResponse, Error error) {
	var newCart entity.Cart
	var cart entity.Cart
	var product entity.Product

	if Error = Repository.Mysql.Where("deleted_at IS NULL").First(&product, Request.ProductId).Error; Error != nil {
		return Response, Error
	}

	if Error = Repository.Mysql.Where("deleted_at IS NULL AND product_id = ?", Request.ProductId).First(&cart, Id).Error; Error != nil {
		return Response, Error
	}

	mapCartRequestToCartEntity(&newCart, Request)
	cart.Qty = newCart.Qty
	cart.UpdateAt = time.Now()

	if Error = Repository.Mysql.Save(&cart).Error; Error != nil {
		return Response, Error
	}

	mapCartEntityToCartResponse(&Response, cart, product)
	return Response, Error
}

func (Repository TransactionRepositoryImpl) DeleteCart(Id int64) (Error error) {
	var cart entity.Cart
	if Error = Repository.Mysql.Where("deleted_at IS NULL").First(&cart, Id).Error; Error != nil {
		return Error
	}

	if Error = Repository.Mysql.Model(&cart).Update("deleted_at", time.Now()).Error; Error != nil {
		return Error
	}
	return Error
}

func (Repository TransactionRepositoryImpl) SubmitTransaction(Request model.TransactionRequest) (Response model.TransactionResponse, Error error) {
	products := make([]entity.Product, len(Request.Details))
	var transaction entity.Transaction
	transactionDetails := make([]entity.TransactionDetail, len(Request.Details))

	var GrandTotal int

	for i, requestProduct := range Request.Details {

		if Error = Repository.Mysql.Where("deleted_at IS NULL").First(&products[i], Request.Details[i].ProductId).Error; Error != nil {
			return Response, Error
		}

		if products[i].Stock < requestProduct.Qty {
			Error = errors.New("STOCK_UNAVAILABLE")
			return Response, Error
		}

		transactionDetails[i] = entity.TransactionDetail{
			Qty:       requestProduct.Qty,
			SubTotal:  requestProduct.Qty * products[i].Price,
			ProductId: requestProduct.ProductId,
		}

		GrandTotal += requestProduct.Qty * products[i].Price
	}

	transaction.Invoice = "ORD" + helper.RandString(5) + strconv.FormatInt(time.Now().Unix(), 10)
	transaction.GrandTotal = GrandTotal
	transaction.AdminId = Request.AdminId

	if Error = Repository.Mysql.Create(&transaction).Error; Error != nil {
		return Response, Error
	}

	for i := range Request.Details {
		transactionDetails[i].TransactionId = transaction.Id

		if Error = Repository.Mysql.Create(&transactionDetails[i]).Error; Error != nil {
			return Response, Error
		}

		products[i].Stock -= transactionDetails[i].Qty
		if Error = Repository.Mysql.Save(&products[i]).Error; Error != nil {
			return Response, Error
		}
	}

	Response.Invoice = transaction.Invoice
	Response.Details = make([]model.TransactionDetailResponse, len(Request.Details))
	for i, detail := range transactionDetails {
		mapTransactionDetailEntityToTransactionDetailResponse(&Response.Details[i], detail, products[i])
	}
	Response.GrandTotal = GrandTotal
	Response.AdminId = Request.AdminId
	Response.AdminName = Request.AdminName
	Response.CreatedAt = transaction.CreatedAt.Format(time.RFC3339)
	Response.UpdatedAt = transaction.UpdatedAt.Format(time.RFC3339)

	return Response, Error
}

func NewTransactionRepository(Mysql *gorm.DB) TransactionRepository {
	return &TransactionRepositoryImpl{
		Mysql: *Mysql,
	}
}
