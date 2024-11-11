package handler

import (
	"love-remittance-be-apps/lib/model"

	"github.com/gofiber/fiber/v2"
)

func handler(f func(c *fiber.Ctx) model.Response) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		response := f(c)
		return c.Status(response.Status).JSON(response)
	}
}

func imageHandler(f func(c *fiber.Ctx) error) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		response := f(c)
		return response
	}
}
