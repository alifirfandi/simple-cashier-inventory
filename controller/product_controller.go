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

func (Controller ProductController) GetAllProducts(c *fiber.Ctx) error {
	query := new(model.ProductRequestQuery)
	if err := c.QueryParser(query); err != nil {
		return exception.ErrorHandler(c, err)
	}
	if query.Sort == "" {
		query.Sort = "id_asc"
	}
	if query.Page <= 0 {
		query.Page = 1
	}

	response, err := Controller.ProductService.GetAllProducts(model.ProductRequestQuery{
		Q:    query.Q,
		Sort: query.Sort,
		Page: query.Page,
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

func (Controller ProductController) GetProductDetail(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return exception.ErrorHandler(c, err)
	}

	if id <= 0 {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Code:   fiber.StatusNotFound,
			Status: "NOT_FOUND",
			Data:   nil,
			Error:  nil,
		})
	}

	response, err := Controller.ProductService.GetProductById(int64(id))
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

func (Controller ProductController) UpdateProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return exception.ErrorHandler(c, err)
	}

	if id <= 0 {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Code:   fiber.StatusNotFound,
			Status: "NOT_FOUND",
			Data:   nil,
			Error:  nil,
		})
	}

	request := new(model.ProductRequest)
	if err = c.BodyParser(request); err != nil {
		return exception.ErrorHandler(c, err)
	}

	file, _ := c.FormFile("image")
	response, err := Controller.ProductService.UpdateProductById(
		int64(id),
		model.ProductRequest{
			Name:  request.Name,
			Price: request.Price,
			Stock: request.Stock,
		},
		file,
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

func (Controller ProductController) DeleteProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return exception.ErrorHandler(c, err)
	}

	if id <= 0 {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Code:   fiber.StatusNotFound,
			Status: "NOT_FOUND",
			Data:   nil,
			Error:  nil,
		})
	}

	err = Controller.ProductService.DeleteProductById(int64(id))
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

func (Controller ProductController) Route(App fiber.Router) {
	router := App.Group("/product")
	router.Post("", middleware.CheckToken(), middleware.CheckRole("SUPERADMIN"), Controller.CreateProduct)
	router.Get("", middleware.CheckToken(), middleware.CheckRole("SUPERADMIN"), Controller.GetAllProducts)
	router.Get("/:id", middleware.CheckToken(), middleware.CheckRole("SUPERADMIN"), Controller.GetProductDetail)
	router.Put("/:id", middleware.CheckToken(), middleware.CheckRole("SUPERADMIN"), Controller.UpdateProduct)
	router.Delete("/:id", middleware.CheckToken(), middleware.CheckRole("SUPERADMIN"), Controller.DeleteProduct)
}

func NewProductController(ProductService *service.ProductService) ProductController {
	return ProductController{
		ProductService: *ProductService,
	}
}
