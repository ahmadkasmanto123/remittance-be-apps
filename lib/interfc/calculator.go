package interfc

import (
	"love-remittance-be-apps/lib/model"

	"github.com/gofiber/fiber/v2"
)

type CalculatorRepository interface {
	GetPayment() ([]model.CalculotorPayment, *model.ErrorData)
	Destination(initiator_id int) ([]model.Destination, *model.ErrorData)
	Transaction(initiator_id int, country_id int) ([]model.Transaction, *model.ErrorData)
	PartnerMasterStoreId(masterStoreId int) (*model.MasterStore, *model.ErrorData)
}

type CalculatorService interface {
	PaymentMethod(c *fiber.Ctx) model.Response
	DestinationPayment(ctx *fiber.Ctx) model.Response
	Transaction(ctx *fiber.Ctx) model.Response
	GetPrice(ctx *fiber.Ctx) model.Response
}
