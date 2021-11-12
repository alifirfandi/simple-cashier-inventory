package controller

import (
	"github.com/alifirfandi/simple-cashier-inventory/exception"
	"github.com/alifirfandi/simple-cashier-inventory/middleware"
	"github.com/alifirfandi/simple-cashier-inventory/model"
	"github.com/alifirfandi/simple-cashier-inventory/service"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	UserService service.UserService
}

func NewUserController(UserService *service.UserService) UserController {
	return UserController{
		UserService: *UserService,
	}
}

func (Controller UserController) Route(App fiber.Router) {
	router := App.Group("/user")
	router.Get("/profile", middleware.CheckToken(), Controller.Profile)
}

func (Controller UserController) Profile(c *fiber.Ctx) error {
	response, err := Controller.UserService.Profile(model.ProfileRequest{
		Email: c.Locals("email").(string),
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
