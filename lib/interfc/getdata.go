package interfc

import (
	"love-remittance-be-apps/lib/model"

	"github.com/gofiber/fiber/v2"
)

type GetDataAllRepository interface {
	GetTransaction(param model.Param, sopId int) ([]model.DataDestination, *model.Pagination, *model.ErrorData)
	GetCountry(param model.Param) ([]model.DataIdName, *model.ErrorData)
	GetProvince(param model.Param) ([]model.DataIdName, *model.Pagination, *model.ErrorData)
	GetCity(param model.Param, provinceId int) ([]model.DataIdName, *model.Pagination, *model.ErrorData)
	GetOccupation() ([]model.DataIdName, *model.ErrorData)
	GetIdentityType() ([]model.DataIdName, *model.ErrorData)
	DataAccountAll(phonePrefix string, phone string) (*model.DataAccountDetail, *model.ErrorData)
	GetSoFund() ([]model.DataIdName, *model.ErrorData)
	GetPurpose() ([]model.DataIdName, *model.ErrorData)
	GetRelations() ([]model.DataIdName, *model.ErrorData)

	GetOccupationIntl(masterStoreId int) ([]model.DataIdName, *model.ErrorData)
	GetIdentityTypeIntl(masterStoreId int) ([]model.DataIdName, *model.ErrorData)
	GetSoFundIntl(masterStoreId int) ([]model.DataIdName, *model.ErrorData)
	GetPurposeIntl(masterStoreId int) ([]model.DataIdName, *model.ErrorData)
	GetRelationsIntl(masterStoreId int) ([]model.DataIdName, *model.ErrorData)

	GetHistory(phone string) ([]model.History, *model.ErrorData)
	GetHistoryDetail(transactionId int) (*model.History, *model.ErrorData)
}

type GetDataAllService interface {
	GetImage(ctx *fiber.Ctx) error
	GetCountry(c *fiber.Ctx) model.Response
	GetProvince(c *fiber.Ctx) model.Response
	GetCity(c *fiber.Ctx) model.Response
	GetOccupation(c *fiber.Ctx) model.Response
	GetIdentityType(c *fiber.Ctx) model.Response
	GetSoFund(c *fiber.Ctx) model.Response
	GetPurpose(c *fiber.Ctx) model.Response
	GetRelations(c *fiber.Ctx) model.Response

	GetOccupationIntl(c *fiber.Ctx) model.Response
	GetIdentityTypeIntl(c *fiber.Ctx) model.Response
	GetSoFundIntl(c *fiber.Ctx) model.Response
	GetPurposeIntl(c *fiber.Ctx) model.Response
	GetRelationsIntl(c *fiber.Ctx) model.Response

	GetAccountDetail(c *fiber.Ctx) model.Response
	GetFromSender(c *fiber.Ctx) model.Response
	GetFromReceiver(c *fiber.Ctx) model.Response
	GetFromAdditional(c *fiber.Ctx) model.Response

	GetHistory(c *fiber.Ctx) model.Response
	GetHistoryDetail(c *fiber.Ctx) model.Response
}
