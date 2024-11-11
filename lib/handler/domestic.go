package handler

import (
	"love-remittance-be-apps/lib/repository"
	"love-remittance-be-apps/lib/service"

	"github.com/gofiber/fiber/v2"
)

func domesticHandler(api fiber.Router) {
	domesticRepository := repository.NewDomesticRepository()
	domesticService := service.NewDomesticService(domesticRepository)

	api.Get("/sop", handler(domesticService.Sopayment))
	api.Get("/destination", handler(domesticService.Destination))
	api.Get("/getPrice", handler(domesticService.GetPrice))
}

func domTrxHandler(api fiber.Router) {
	domesticRepository := repository.NewDomesticRepository()
	domesticService := service.NewDomesticService(domesticRepository)

	api.Get("/checkAccount", handler(domesticService.CheckAccount))
	api.Get("/createPaycode", handler(domesticService.CreatePaycode))
}
