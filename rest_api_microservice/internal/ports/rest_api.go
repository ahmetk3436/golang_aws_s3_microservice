package ports

import "github.com/gofiber/fiber/v3"

type RestApi interface {
	GetProductById(ctx *fiber.Ctx) error
	Start(port int) error
	Stop() error
}
