package service

import (
	"github.com/alifirfandi/simple-cashier-inventory/model"
	"github.com/alifirfandi/simple-cashier-inventory/repository"
	"github.com/alifirfandi/simple-cashier-inventory/validation"
)

type TransactionServiceImpl struct {
	TransactionRepository repository.TransactionRepository
}

func (Service TransactionServiceImpl) GetUserCart(idUser int64) (Response model.CartListResponse, Error error) {
	Response, Error = Service.TransactionRepository.GetAllCart(idUser)

	return Response, Error
}

func (Service TransactionServiceImpl) InsertCart(Request model.CartRequest) (Response model.CartResponse, Error error) {
	if Error = validation.CartValidation(Request); Error != nil {
		return Response, Error
	}

	Response, Error = Service.TransactionRepository.InsertCart(Request)
	return Response, Error
}

func (Service TransactionServiceImpl) UpdateCart(Id int64, Request model.CartRequest) (Response model.CartResponse, Error error) {
	if Error = validation.CartValidation(Request); Error != nil {
		return Response, Error
	}

	Response, Error = Service.TransactionRepository.UpdateCart(Id, Request)
	return Response, Error
}

func (Service TransactionServiceImpl) DeleteCart(Id int64) (Error error) {
	if Error = Service.TransactionRepository.DeleteCart(Id); Error != nil {
		return Error
	}

	return nil
}

func (Service TransactionServiceImpl) SubmitTransaction(Request model.TransactionRequest) (Response model.TransactionResponse, Error error) {
	for _, detail := range Request.Details {
		if Error = validation.CartValidation(detail); Error != nil {
			return Response, Error
		}
	}

	Response, Error = Service.TransactionRepository.SubmitTransaction(Request)
	return Response, Error
}

func NewTransactionService(TransactionRepository *repository.TransactionRepository) TransactionService {
	return &TransactionServiceImpl{
		TransactionRepository: *TransactionRepository,
	}
}
