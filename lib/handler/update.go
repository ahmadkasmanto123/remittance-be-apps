package handler

import (
	"love-remittance-be-apps/lib/repository"
	"love-remittance-be-apps/lib/service"

	"github.com/gofiber/fiber/v2"
)

func updateHandler(api fiber.Router) {
	newRepository := repository.NewUpdateRepository()
	newService := service.NewUpdateService(newRepository)

	api.Post("/profile", handler(newService.Profile))
}
