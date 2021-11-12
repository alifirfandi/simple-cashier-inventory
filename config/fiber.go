package config

import (
	"github.com/alifirfandi/simple-cashier-inventory/exception"

	"github.com/gofiber/fiber/v2"
)

func NewFiberConfig() fiber.Config {
	return fiber.Config{
		ErrorHandler: exception.ErrorHandler,
	}
}
