package interfc

import (
	"love-remittance-be-apps/lib/model"

	"github.com/gofiber/fiber/v2"
)

type RegistrationRepository interface {
	CountAccount(phone string, phonePrefix string) (*model.CountSomething, *model.ErrorData)
	CountEmail(email string) (*model.CountSomething, *model.ErrorData)
	CountMsAccount(phone string) (*model.CountSomething, *model.ErrorData)
	Storecheck(phonePrefix string) (*model.StoreCheck, *model.ErrorData)
	UpdateMsAccount(userName string, phone string) (*int, *model.ErrorData)
	InsertMsAccount(userName string, phone string) (*int, *model.ErrorData)
	InsertMsStoreAdmin(req model.DefaultRequest[model.CreateCustomer], storeId int, accountId int) (*int, *model.ErrorData)
	DataAccountReg(phone string, phonePrefix string) (*model.DataAccount, *model.ErrorData)
	UpdateMsStoreAdmin(deviceId string, adminId int) *model.ErrorData
}

type RegistrationService interface {
	NewCustomer(ctx *fiber.Ctx) model.Response
	RegOtpValidate(ctx *fiber.Ctx) model.Response
}
