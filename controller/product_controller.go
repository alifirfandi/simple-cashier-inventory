package controller

import (
	"github.com/alifirfandi/simple-cashier-inventory/exception"
	"github.com/alifirfandi/simple-cashier-inventory/middleware"
	"github.com/alifirfandi/simple-cashier-inventory/model"
	"github.com/alifirfandi/simple-cashier-inventory/service"
	"github.com/gofiber/fiber/v2"
)

type ProductController struct {
	ProductService service.ProductService
}

func NewProductController(ProductService *service.ProductService) ProductController {
	return ProductController{
		ProductService: *ProductService,
	}
}

func (Controller ProductController) Route(App fiber.Router) {
	router := App.Group("/product")
	router.Post("", middleware.CheckToken(), middleware.CheckRole("SUPERADMIN"), Controller.CreateProduct)
}

func (Controller ProductController) CreateProduct(c *fiber.Ctx) error {
	request := new(model.ProductRequest)
	if err := c.BodyParser(request); err != nil {
		return exception.ErrorHandler(c, err)
	}

	file, _ := c.FormFile("image")
	response, err := Controller.ProductService.InsertProduct(model.ProductRequest{
		Name:  request.Name,
		Price: request.Price,
		Stock: request.Stock,
	}, file)
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
