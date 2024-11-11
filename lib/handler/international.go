package handler

import (
	"love-remittance-be-apps/lib/repository"
	"love-remittance-be-apps/lib/service"

	"github.com/gofiber/fiber/v2"
)

func internationalHandler(api fiber.Router) {
	newRepository := repository.NewInternationalRepository()
	newService := service.NewInternationalService(newRepository)

	api.Get("/sop", handler(newService.SourceOfPayment))
	api.Get("/country", handler(newService.AvailableCountry))
	api.Get("/destination", handler(newService.DestinationByCountry))
	api.Get("/getPrice", handler(newService.GetPrice))
}

func intlTrxHandler(api fiber.Router) {
	newRepository := repository.NewInternationalRepository()
	newService := service.NewInternationalService(newRepository)

	api.Post("/createPaycode", handler(newService.CreatePaycode))
}
