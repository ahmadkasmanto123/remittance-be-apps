package config

import (
	"log"
	"love-remittance-be-apps/lib/model"

	"github.com/gofiber/fiber/v2"
)

func errorHandlerConfig(ctx *fiber.Ctx, err error) error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	return ctx.Status(fiber.StatusInternalServerError).JSON(model.Response{
		RC:      "1006",
		Message: "Internal Server Error",
		Errors: append([]model.ErrorData{}, model.ErrorData{
			Description: err.Error(),
		}),
	})
}
