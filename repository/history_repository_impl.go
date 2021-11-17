package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/alifirfandi/simple-cashier-inventory/entity"
	"github.com/alifirfandi/simple-cashier-inventory/model"
	"gorm.io/gorm"
)

func mapTrxDetailEntityToHistoryDetailResponse(historyDetailResponse *model.TransactionDetailResponse, trxEntityDetail entity.TransactionDetail) {
	historyDetailResponse.Id = trxEntityDetail.Id
	historyDetailResponse.ProductId = trxEntityDetail.ProductId
	historyDetailResponse.Name = trxEntityDetail.Product.Name
	historyDetailResponse.ImageUrl = trxEntityDetail.Product.ImageUrl
	historyDetailResponse.Price = trxEntityDetail.Product.Price
	historyDetailResponse.Qty = trxEntityDetail.Qty
	historyDetailResponse.SubTotal = trxEntityDetail.SubTotal
}

func mapTrxEntityToHistoryResponse(historyResponse *model.HistoryResponse, transactionEntity entity.Transaction) {
	historyResponse.Invoice = transactionEntity.Invoice
	historyResponse.Details = make([]model.TransactionDetailResponse, len(transactionEntity.TransactionDetails))
	for i, detail := range transactionEntity.TransactionDetails {
		mapTrxDetailEntityToHistoryDetailResponse(&historyResponse.Details[i], detail)
	}
	historyResponse.GrandTotal = transactionEntity.GrandTotal
	historyResponse.AdminId = transactionEntity.AdminId
	historyResponse.AdminName = transactionEntity.Admin.Name
	historyResponse.CreatedAt = transactionEntity.CreatedAt.Format(time.RFC3339)
	historyResponse.UpdatedAt = transactionEntity.UpdatedAt.Format(time.RFC3339)
}

type HistoryRepositoryImpl struct {
	Mysql gorm.DB
}

func (Repository HistoryRepositoryImpl) GetAllHistories(Query model.HistorySelectQuery) (Response []model.HistoryResponse, Error error) {
	var transactions []entity.Transaction

	var q strings.Builder
	q.WriteString("deleted_at IS NULL AND invoice LIKE ?")
	if Query.StartDate != "" {
		q.WriteString(fmt.Sprintf(" AND created_at >= '%s'", Query.StartDate))
	}
	if Query.EndDate != "" {
		q.WriteString(fmt.Sprintf(" AND created_at <= '%s'", Query.EndDate))
	}

	Error = Repository.Mysql.Where(q.String(), Query.Search).Limit(Query.Limit).Offset(Query.Start).Preload("TransactionDetails.Product").Preload("Admin").Find(&transactions).Error
	if Error != nil {
		return Response, Error
	}

	Response = make([]model.HistoryResponse, len(transactions))
	for i, transaction := range transactions {
		mapTrxEntityToHistoryResponse(&Response[i], transaction)
	}
	return Response, Error
}

func (Repository HistoryRepositoryImpl) GetHistoryByInvoice(Invoice string) (Response model.HistoryResponse, Error error) {
	var transaction entity.Transaction
	Error = Repository.Mysql.Where("invoice = ? AND deleted_at IS NULL", Invoice).Preload("TransactionDetails.Product").Preload("Admin").First(&transaction).Error
	if Error != nil {
		return Response, Error
	}
	mapTrxEntityToHistoryResponse(&Response, transaction)
	return
}

func (Repository HistoryRepositoryImpl) CountHistories() (Result int64, Error error) {
	Error = Repository.Mysql.Model(&entity.Transaction{}).Where("deleted_at IS NULL").Count(&Result).Error
	return Result, Error
}

func NewHistoryRepository(Mysql *gorm.DB) HistoryRepository {
	return &HistoryRepositoryImpl{
		Mysql: *Mysql,
	}
}
