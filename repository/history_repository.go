package repository

import "github.com/alifirfandi/simple-cashier-inventory/model"

type HistoryRepository interface {
	GetAllHistories() (Response []model.HistoryResponse, Error error)
	GetHistoryByInvoice(Invoice string) (Response model.HistoryResponse, Error error)
}
