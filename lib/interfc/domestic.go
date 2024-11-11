package interfc

import (
	"love-remittance-be-apps/lib/model"

	"github.com/gofiber/fiber/v2"
)

type DomesticRepository interface {
	GetSourceOfPayment() ([]model.DataSOPayment, *model.ErrorData)
	GetDestination(param model.Param, sopId int) ([]model.DataDestination, *model.Pagination, *model.ErrorData)
	PartnerMasterStoreId(masterStoreId int) (*model.MasterStore, *model.ErrorData)
	DataInitiator(sopId int) (*model.DataSOPayment, *model.ErrorData)
	DataAccountAll(phonePrefix string, phone string) (*model.DataAccountDetail, *model.ErrorData)
}

type DomesticService interface {
	Sopayment(c *fiber.Ctx) model.Response
	Destination(c *fiber.Ctx) model.Response
	GetPrice(ctx *fiber.Ctx) model.Response
	CheckAccount(ctx *fiber.Ctx) model.Response
	CreatePaycode(ctx *fiber.Ctx) model.Response
}
