package service

import (
	"fmt"
	"math"
	"os"
	"strconv"

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
	if Query.Page <= 0 {
		Query.Page = 1
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

func NewHistoryService(Repository *repository.HistoryRepository) HistoryService {
	return &HistoryServiceImpl{
		Repository: *Repository,
	}
}
