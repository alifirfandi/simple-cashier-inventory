package repository

import (
	"github.com/alifirfandi/simple-cashier-inventory/model"
	"gorm.io/gorm"
)

type HistoryRepositoryImpl struct {
	Mysql gorm.DB
}

func (Repository HistoryRepositoryImpl) GetAllHistories() (Response []model.HistoryResponse, Error error) {
	return
}

func (Repository HistoryRepositoryImpl) GetHistoryByInvoice(Invoice string) (Response model.HistoryResponse, Error error) {
	return
}

func NewHistoryRepository(Mysql *gorm.DB) HistoryRepository {
	return &HistoryRepositoryImpl{
		Mysql: *Mysql,
	}
}
