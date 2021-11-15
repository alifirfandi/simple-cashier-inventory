package controller

import (
	"github.com/alifirfandi/simple-cashier-inventory/middleware"
	"github.com/alifirfandi/simple-cashier-inventory/service"
	"github.com/gofiber/fiber/v2"
)

type HistoryController struct {
	Service service.HistoryService
}

func (Controller HistoryController) GetHistoryList(c *fiber.Ctx) error {
	err := c.SendString("Hello")
	return err
}

func (Controller HistoryController) GetHistoryDetail(c *fiber.Ctx) error {
	err := c.SendString("Hi")
	return err
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
