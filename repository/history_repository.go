package repository

import "github.com/alifirfandi/simple-cashier-inventory/model"

type HistoryRepository interface {
	GetAllHistories(Query model.HistorySelectQuery) (Response []model.HistoryResponse, Error error)
	GetHistoryByInvoice(Invoice string) (Response model.HistoryResponse, Error error)
	CountHistories() (Result int64, Error error)
}
