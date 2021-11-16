package main

import (
	"os"

	"github.com/alifirfandi/simple-cashier-inventory/config"
	"github.com/alifirfandi/simple-cashier-inventory/controller"
	"github.com/alifirfandi/simple-cashier-inventory/exception"
	"github.com/alifirfandi/simple-cashier-inventory/repository"
	"github.com/alifirfandi/simple-cashier-inventory/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	mysql := config.MysqlConnection()

	authRepository := repository.NewAuthRepository(mysql)
	userRepository := repository.NewUserRepository(mysql)
	productRepository := repository.NewProductRepository(mysql)
	historyRepository := repository.NewHistoryRepository(mysql)

	authService := service.NewAuthService(&authRepository)
	userService := service.NewUserService(&userRepository)
	productService := service.NewProductService(&productRepository)
	historyService := service.NewHistoryService(&historyRepository)

	authController := controller.NewAuthController(&authService)
	userController := controller.NewUserController(&userService)
	productController := controller.NewProductController(&productService)
	historyController := controller.NewHistoryController(&historyService)

	app := fiber.New(config.NewFiberConfig())
	app.Use(recover.New())

	v1 := app.Group("/api/v1")
	authController.Route(v1)
	userController.Route(v1)
	productController.Route(v1)
	historyController.Route(v1)

	// Start App
	err := app.Listen(os.Getenv("PORT"))
	exception.PanicIfNeeded(err)
}
