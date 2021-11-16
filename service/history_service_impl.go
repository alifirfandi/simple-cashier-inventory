package service

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/alifirfandi/simple-cashier-inventory/helper"

	"github.com/alifirfandi/simple-cashier-inventory/model"
	"github.com/alifirfandi/simple-cashier-inventory/repository"
)

type HistoryServiceImpl struct {
	Repository repository.HistoryRepository
}

func (Service HistoryServiceImpl) GetAllHistories(Query model.HistoryRequestQuery) (Response model.HistoryListResponse, Error error) {
	limit, Error := strconv.Atoi(os.Getenv("LIMIT_PER_PAGE"))
	if Error != nil {
		return Response, Error
	}

	page := (Query.Page - 1) * limit
	histories, Error := Service.Repository.GetAllHistories(model.HistorySelectQuery{
		Search:    fmt.Sprintf("%%%s%%", Query.Q),
		Start:     page,
		Limit:     limit,
		StartDate: Query.StartDate,
		EndDate:   Query.EndDate,
	})
	if Error != nil {
		return Response, Error
	}
	historiesCount, Error := Service.Repository.CountHistories()
	if Error != nil {
		return Response, Error
	}

	Response = model.HistoryListResponse{
		TotalData:    int(historiesCount),
		TotalPage:    int(math.Ceil(float64(historiesCount) / float64(limit))),
		CurrentPage:  Query.Page,
		LimitPerPage: limit,
		Histories:    histories,
	}
	return Response, Error
}

func (Service HistoryServiceImpl) GetHistoryByInvoice(Invoice string) (Response model.HistoryResponse, Error error) {
	Response, Error = Service.Repository.GetHistoryByInvoice(Invoice)
	return Response, Error
}

func (Service HistoryServiceImpl) ExportPdfByInvoice(Invoice string) (Response string, Error error) {
	res, Error := Service.Repository.GetHistoryByInvoice(Invoice)
	if Error != nil {
		return "", Error
	}

	invoiceData := helper.InvoiceData{
		Invoice:    res.Invoice,
		Name:       res.AdminName,
		Date:       res.CreatedAt,
		GrandTotal: res.GrandTotal,
		Items:      make([]helper.InvoiceItemData, len(res.Details)),
	}
	for i, detail := range res.Details {
		invoiceData.Items[i] = helper.InvoiceItemData{
			Name:     detail.Name,
			Price:    detail.Price,
			Qty:      detail.Qty,
			SubTotal: detail.SubTotal,
		}
	}

	fileName := fmt.Sprintf("%s-%d-%s.pdf", helper.RandString(9), time.Now().Unix(), res.Invoice)
	Error = helper.InvoiceExporter(helper.InvoiceTemplate, os.Getenv("PDF_LOCATION"), fileName, invoiceData)
	if Error != nil {
		return "", Error
	}

	return "", Error
}

func NewHistoryService(Repository *repository.HistoryRepository) HistoryService {
	return &HistoryServiceImpl{
		Repository: *Repository,
	}
}
