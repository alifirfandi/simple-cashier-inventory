package service

import (
	"github.com/alifirfandi/simple-cashier-inventory/model"
)

type TransactionService interface {
	InsertCart(Request model.CartRequest) (Response model.CartResponse, Error error)
	GetUserCart(idUser int64) (Response model.CartListResponse, Error error)
	UpdateCart(Id int64, Request model.CartRequest) (Response model.CartResponse, Error error)
	DeleteCart(Id int64) (Error error)
	SubmitTransaction(Request model.TransactionRequest) (Response model.TransactionResponse, Error error)
}
