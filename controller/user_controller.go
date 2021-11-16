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

func (Controller UserController) CreateUser(c *fiber.Ctx) error {
	request := new(model.UserRequest)
	if err := c.BodyParser(request); err != nil {
		return exception.ErrorHandler(c, err)
	}

	response, err := Controller.UserService.InsertUser(model.UserRequest{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
		Role:     request.Role,
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

func (Controller UserController) GetAllUser(c *fiber.Ctx) error {
	QueryParams := new(model.UserSelectQuery)

	if err := c.QueryParser(QueryParams); err != nil {
		return err
	}

	response, err := Controller.UserService.GetAllUser(*QueryParams)
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

func (Controller UserController) GetUserDetail(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return exception.ErrorHandler(c, err)
	}

	if id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Code:   fiber.StatusBadRequest,
			Status: "BAD_REQUEST",
			Data:   nil,
			Error: map[string]string{
				"id": "INVALID_ID",
			},
		})
	}

	response, err := Controller.UserService.GetUserById(int64(id))
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

func (Controller UserController) UpdateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return exception.ErrorHandler(c, err)
	}

	if id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Code:   fiber.StatusBadRequest,
			Status: "BAD_REQUEST",
			Data:   nil,
			Error: map[string]string{
				"id": "INVALID_ID",
			},
		})
	}

	request := new(model.UserRequest)
	if err = c.BodyParser(request); err != nil {
		return exception.ErrorHandler(c, err)
	}

	response, err := Controller.UserService.UpdateUser(
		int64(id),
		model.UserRequest{
			Name:     request.Name,
			Email:    request.Email,
			Password: request.Password,
			Role:     request.Role,
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

func (Controller UserController) DeleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return exception.ErrorHandler(c, err)
	}

	if id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Code:   fiber.StatusBadRequest,
			Status: "BAD_REQUEST",
			Data:   nil,
			Error: map[string]string{
				"id": "INVALID_ID",
			},
		})
	}

	err = Controller.UserService.DeleteUser(int64(id))
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

func (Controller UserController) Route(App fiber.Router) {
	router := App.Group("/admin")
	router.Post("", middleware.CheckToken(), middleware.CheckRole("SUPERADMIN"), Controller.CreateUser)
	router.Get("", middleware.CheckToken(), middleware.CheckRole("SUPERADMIN"), Controller.GetAllUser)
	router.Get("/:id", middleware.CheckToken(), middleware.CheckRole("SUPERADMIN"), Controller.GetUserDetail)
	router.Put("/:id", middleware.CheckToken(), middleware.CheckRole("SUPERADMIN"), Controller.UpdateUser)
	router.Delete("/:id", middleware.CheckToken(), middleware.CheckRole("SUPERADMIN"), Controller.DeleteUser)
}
