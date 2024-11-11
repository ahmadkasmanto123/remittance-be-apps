package handler

import (
	"love-remittance-be-apps/lib/repository"
	"love-remittance-be-apps/lib/service"

	"github.com/gofiber/fiber/v2"
)

func getDataHandler(api fiber.Router) {
	newRepository := repository.NewGetDataAllRepository()
	newService := service.NewGetDataAllService(newRepository)

	api.Get("/image", imageHandler(newService.GetImage))

	api.Get("/country", handler(newService.GetCountry))
	api.Get("/province", handler(newService.GetProvince))
	api.Get("/city", handler(newService.GetCity))
	api.Get("/occupation", handler(newService.GetOccupation))
	api.Get("/identityType", handler(newService.GetIdentityType))
	api.Get("/sofunding", handler(newService.GetSoFund))
	api.Get("/purpose", handler(newService.GetPurpose))
	api.Get("/relationship", handler(newService.GetRelations))

	api.Get("/account", handler(newService.GetAccountDetail))

	api.Get("/form/sender", handler(newService.GetFromSender))
	api.Get("/form/receiver", handler(newService.GetFromReceiver))
	api.Get("/form/additional", handler(newService.GetFromAdditional))

	api.Get("/history", handler(newService.GetHistory))
	api.Get("/history/detail", handler(newService.GetHistoryDetail))

	api.Get("/occupation/intl", handler(newService.GetOccupationIntl))
	api.Get("/identityType/intl", handler(newService.GetIdentityTypeIntl))
	api.Get("/sofunding/intl", handler(newService.GetSoFundIntl))
	api.Get("/purpose/intl", handler(newService.GetPurposeIntl))
	api.Get("/relationship/intl", handler(newService.GetRelationsIntl))
}
