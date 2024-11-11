package handler

import (
	"love-remittance-be-apps/lib/repository"
	"love-remittance-be-apps/lib/service"

	"github.com/gofiber/fiber/v2"
)

func calculatorHandler(api fiber.Router) {
	calculatorRepository := repository.NewCalculatorRepository()
	calculatorService := service.NewCalculatorService(calculatorRepository)

	api.Get("/paymentMethod", handler(calculatorService.PaymentMethod))
	api.Get("/destination", handler(calculatorService.DestinationPayment))
	api.Get("/transaction", handler(calculatorService.Transaction))
	api.Get("/price", handler(calculatorService.GetPrice))

}
