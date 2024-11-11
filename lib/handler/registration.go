package handler

import (
	"love-remittance-be-apps/lib/repository"
	"love-remittance-be-apps/lib/service"

	"github.com/gofiber/fiber/v2"
)

func registerHandler(api fiber.Router) {
	newRepository := repository.NewRegistrationRepository()
	newService := service.NewRegistrationService(newRepository)

	api.Post("/create", handler(newService.NewCustomer))
	api.Post("/validated", handler(newService.RegOtpValidate))
}
