package service

import (
	"github.com/alifirfandi/simple-cashier-inventory/model"
	"github.com/alifirfandi/simple-cashier-inventory/repository"
)

type HistoryServiceImpl struct {
	Repository repository.HistoryRepository
}

func (Service HistoryServiceImpl) GetAllHistories() (Response model.HistoryListResponse, Error error) {
	return
}

func (Service HistoryServiceImpl) GetHistoryByInvoice(Invoice string) (Response model.HistoryResponse, Error error) {
	return
}

func NewHistoryService(Repository *repository.HistoryRepository) HistoryService {
	return &HistoryServiceImpl{
		Repository: *Repository,
	}
}
