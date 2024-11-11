package interfc

import (
	"love-remittance-be-apps/lib/model"

	"github.com/gofiber/fiber/v2"
)

type InternationalRepository interface {
	GetSourceOfPayment() ([]model.DataSOPayment, *model.ErrorData)
	GetCountry(sopId int) ([]model.CountryData, *model.ErrorData)
	GetDestination(param model.Param, countryId int, sopiId int) ([]model.BankDest, *model.ErrorData)
	GetCredential(masterStoreId int, sopId int) (*model.Credential, *model.ErrorData)
	DataAccountAll(phonePrefix string, phone string) (*model.DataAccountDetail, *model.ErrorData)
	GetIdentitas(id int) (*model.DataIdNkey, *model.ErrorData)
	GetRelations(id int) (*model.DataIdNkey, *model.ErrorData)
	GetPurposes(id int) (*model.DataIdNkey, *model.ErrorData)
	GetFunding(id int) (*model.DataIdNkey, *model.ErrorData)
	GetOccupation(id int) (*model.DataIdNkey, *model.ErrorData)
}

type InternationalService interface {
	SourceOfPayment(c *fiber.Ctx) model.Response
	AvailableCountry(c *fiber.Ctx) model.Response
	DestinationByCountry(c *fiber.Ctx) model.Response
	GetPrice(c *fiber.Ctx) model.Response

	CreatePaycode(c *fiber.Ctx) model.Response
}
