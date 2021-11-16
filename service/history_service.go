package service

import "github.com/alifirfandi/simple-cashier-inventory/model"

type HistoryService interface {
	GetAllHistories(Query model.HistoryRequestQuery) (Response model.HistoryListResponse, Error error)
	GetHistoryByInvoice(Invoice string) (Response model.HistoryResponse, Error error)
	ExportPdfByInvoice(Invoice string) (Response string, Error error)
}
