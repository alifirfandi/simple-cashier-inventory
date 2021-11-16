package controller

import (
	"github.com/alifirfandi/simple-cashier-inventory/exception"
	"github.com/alifirfandi/simple-cashier-inventory/middleware"
	"github.com/alifirfandi/simple-cashier-inventory/model"
	"github.com/alifirfandi/simple-cashier-inventory/service"
	"github.com/gofiber/fiber/v2"
)

type TransactionController struct {
	TransactionService service.TransactionService
}

func (Controller TransactionController) GetUserCart(c *fiber.Ctx) error {
	adminId := c.Locals("id").(uint64)

	response, err := Controller.TransactionService.GetUserCart(int64(adminId))
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

func (Controller TransactionController) CreateCart(c *fiber.Ctx) error {
	request := new(model.CartRequest)
	if err := c.BodyParser(request); err != nil {
		return exception.ErrorHandler(c, err)
	}

	adminId := c.Locals("id").(uint64)
	response, err := Controller.TransactionService.InsertCart(model.CartRequest{
		AdminId:   int64(adminId),
		Qty:       request.Qty,
		ProductId: request.ProductId,
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

func (Controller TransactionController) UpdateCart(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return exception.ErrorHandler(c, err)
	}

	if id <= 0 {
		return c.Status(fiber.StatusOK).JSON(model.Response{
			Code:   fiber.StatusOK,
			Status: "BAD_REQUEST",
			Data:   nil,
			Error: map[string]string{
				"id": "INVALID_ID",
			},
		})
	}

	request := new(model.CartRequest)
	if err = c.BodyParser(request); err != nil {
		return exception.ErrorHandler(c, err)
	}

	adminId := c.Locals("id").(uint64)
	response, err := Controller.TransactionService.UpdateCart(
		int64(id),
		model.CartRequest{
			AdminId:   int64(adminId),
			Qty:       request.Qty,
			ProductId: request.ProductId,
		},
	)
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

func (Controller TransactionController) DeleteCart(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return exception.ErrorHandler(c, err)
	}

	if id <= 0 {
		return c.Status(fiber.StatusOK).JSON(model.Response{
			Code:   fiber.StatusOK,
			Status: "BAD_REQUEST",
			Data:   nil,
			Error: map[string]string{
				"id": "INVALID_ID",
			},
		})
	}

	err = Controller.TransactionService.DeleteCart(int64(id))
	if err != nil {
		return exception.ErrorHandler(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(model.Response{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   nil,
		Error:  nil,
	})
}

func (Controller TransactionController) SubmitTransaction(c *fiber.Ctx) error {
	request := new(model.TransactionRequest)
	if err := c.BodyParser(request); err != nil {
		return exception.ErrorHandler(c, err)
	}

	adminId := c.Locals("id").(uint64)
	adminName := c.Locals("name").(string)
	response, err := Controller.TransactionService.SubmitTransaction(model.TransactionRequest{
		AdminId:   int64(adminId),
		Details:   request.Details,
		AdminName: adminName,
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

func (Controller TransactionController) Route(App fiber.Router) {
	router := App.Group("/transaction")
	router.Get("/cart", middleware.CheckToken(), middleware.CheckRole("ADMIN,SUPERADMIN"), Controller.GetUserCart)
	router.Post("/cart", middleware.CheckToken(), middleware.CheckRole("ADMIN,SUPERADMIN"), Controller.CreateCart)
	router.Put("/cart/:id", middleware.CheckToken(), middleware.CheckRole("ADMIN,SUPERADMIN"), Controller.UpdateCart)
	router.Delete("/cart/:id", middleware.CheckToken(), middleware.CheckRole("ADMIN,SUPERADMIN"), Controller.DeleteCart)
	router.Post("/submit", middleware.CheckToken(), middleware.CheckRole("ADMIN,SUPERADMIN"), Controller.SubmitTransaction)
}

func NewTransactionController(TransactionService *service.TransactionService) TransactionController {
	return TransactionController{
		TransactionService: *TransactionService,
	}
}
