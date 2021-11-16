package controller

import (
	"fmt"

	"github.com/alifirfandi/simple-cashier-inventory/exception"
	"github.com/alifirfandi/simple-cashier-inventory/middleware"
	"github.com/alifirfandi/simple-cashier-inventory/model"
	"github.com/alifirfandi/simple-cashier-inventory/service"
	"github.com/gofiber/fiber/v2"
)

type HistoryController struct {
	Service service.HistoryService
}

func (Controller HistoryController) GetHistoryList(c *fiber.Ctx) error {
	query := new(model.HistoryRequestQuery)
	if err := c.QueryParser(query); err != nil {
		return exception.ErrorHandler(c, err)
	}
	if query.StartDate != "" {
		query.StartDate = fmt.Sprintf("%sT00:00:00Z", query.StartDate)
	}
	if query.EndDate != "" {
		query.EndDate = fmt.Sprintf("%sT00:00:00Z", query.EndDate)
	}

	response, err := Controller.Service.GetAllHistories(model.HistoryRequestQuery{
		Q:         query.Q,
		Page:      query.Page,
		StartDate: query.StartDate,
		EndDate:   query.EndDate,
	})
	if err != nil {
		return exception.ErrorHandler(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(model.Response{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   response,
		Error:  nil,
	})
}

func (Controller HistoryController) GetHistoryDetail(c *fiber.Ctx) error {
	invoice := c.Params("invoice")
	response, err := Controller.Service.GetHistoryByInvoice(invoice)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(model.Response{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   response,
		Error:  nil,
	})
}

func (Controller HistoryController) Route(App fiber.Router) {
	router := App.Group("/history")
	router.Get("", middleware.CheckToken(), middleware.CheckRole("ADMIN,SUPERADMIN"), Controller.GetHistoryList)
	router.Get("/:invoice", middleware.CheckToken(), middleware.CheckRole("ADMIN,SUPERADMIN"), Controller.GetHistoryDetail)
}

func NewHistoryController(Service *service.HistoryService) HistoryController {
	return HistoryController{
		Service: *Service,
	}
}
