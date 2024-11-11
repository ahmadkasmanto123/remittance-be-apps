package interfc

import (
	"love-remittance-be-apps/lib/model"

	"github.com/gofiber/fiber/v2"
)

type LoginRepository interface {
	CountAccount(phone string, phonePrefix string) (*model.CountSomething, *model.ErrorData)
	DataAccount(phone string, phonePrefix string) (*model.DataAccount, *model.ErrorData)
	UpdateAccountStatus(phone string, statusId int) *model.ErrorData
	UpdateMsStoreAdmin(deviceId string, adminId int) *model.ErrorData
}

type LoginService interface {
	Login(ctx *fiber.Ctx) model.Response
	LogInOtpCreate(ctx *fiber.Ctx) model.Response
	LogInOtpValidate(ctx *fiber.Ctx) model.Response
}
