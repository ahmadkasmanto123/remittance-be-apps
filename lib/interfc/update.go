package interfc

import (
	"love-remittance-be-apps/lib/model"

	"github.com/gofiber/fiber/v2"
)

type UpdateRepository interface {
	CountAccount(phone string, phonePrefix string) (*model.CountSomething, *model.ErrorData)
	UpdateMsAccount(req model.DefaultRequest[model.UpdateProfil], img []string) (*int, *model.ErrorData)
	GetDataaaa() ([]model.DataSOPayment, *model.ErrorData)
}

type UpdateService interface {
	Profile(c *fiber.Ctx) model.Response
}
