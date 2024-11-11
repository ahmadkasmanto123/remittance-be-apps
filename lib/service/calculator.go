package service

import (
	"love-remittance-be-apps/core/rc"
	"love-remittance-be-apps/core/utils"
	"love-remittance-be-apps/lib/interfc"
	"love-remittance-be-apps/lib/model"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CalculatorServiceImpl struct {
	CalculatorRepository interfc.CalculatorRepository
}

func NewCalculatorService(repository interfc.CalculatorRepository) interfc.CalculatorService {
	return &CalculatorServiceImpl{
		CalculatorRepository: repository,
	}
}

type PaymentMethodStandard struct {
	Extref string `validate:"required" json:"extref"`
	Lang   string `validate:"required" json:"lang"`
}

func (serv *CalculatorServiceImpl) PaymentMethod(c *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[PaymentMethodStandard](c.Body())

	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "error sini",
		}
	}

	// calculatorRepository := repository.CalculatorRepository{}
	payments, error := serv.CalculatorRepository.GetPayment()

	if error != nil {
		return model.Response{
			Status: error.Status,
			RC:     error.RC,
			Errors: append([]model.ErrorData{}, *error),
		}
	}

	return model.Response{
		Status:  fiber.StatusOK,
		RC:      rc.SUCCESS.String(),
		Message: rc.SUCCESS.Message(),
		Extref:  req.Extref,
		Data:    payments,
	}

}

type DestinationPayment struct {
	Lang        string `validate:"required" json:"lang"`
	Extref      string `validate:"required" json:"extref"`
	InitiatorId string `validate:"required" json:"initiator_id"`
}

func (serv *CalculatorServiceImpl) DestinationPayment(ctx *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[DestinationPayment](ctx.Body())
	if errData != nil {
		return model.Response{
			Status: fiber.StatusBadRequest,
			RC:     rc.FAILED.String(),
			Errors: errData,
		}
	}

	// calculatorRepository := repository.CalculatorRepository{}
	i, errAt := strconv.Atoi(req.InitiatorId)
	if errAt != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Message: "Error Parsing",
		}
	}
	destination, err := serv.CalculatorRepository.Destination(i)

	if err != nil {
		return model.Response{
			Status: err.Status,
			RC:     err.RC,
			Errors: append([]model.ErrorData{}, *err),
		}
	}

	return model.Response{
		Status:  fiber.StatusOK,
		RC:      rc.SUCCESS.String(),
		Message: rc.SUCCESS.Message(),
		Extref:  req.Extref,
		Data:    destination,
	}
}

type TransactionType struct {
	Lang        string `validate:"required" json:"lang"`
	Extref      string `validate:"required" json:"extref"`
	InitiatorId string `validate:"required" json:"initiator_id"`
	CountryId   string `validate:"required" json:"country_id"`
}

func (serv *CalculatorServiceImpl) Transaction(ctx *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[TransactionType](ctx.Body())
	if errData != nil {
		return model.Response{
			Status: fiber.StatusBadRequest,
			RC:     rc.FAILED.String(),
			Errors: errData,
		}
	}

	// calculatorRepository := repository.CalculatorRepository{}
	i, errAt := strconv.Atoi(req.InitiatorId)
	if errAt != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Message: "Error Parsing",
		}
	}
	j, errAt := strconv.Atoi(req.CountryId)
	if errAt != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Message: "Error Parsing",
		}
	}
	transaction, err := serv.CalculatorRepository.Transaction(i, j)

	if err != nil {
		return model.Response{
			Status: err.Status,
			RC:     err.RC,
			Errors: append([]model.ErrorData{}, *err),
		}
	}

	return model.Response{
		Status:  fiber.StatusOK,
		RC:      rc.SUCCESS.String(),
		Message: rc.SUCCESS.Message(),
		Extref:  req.Extref,
		Data:    transaction,
	}
}

type Price struct {
	Reverse           string `validate:"required" json:"reverse"`
	TransactionAmount string `validate:"required" json:"transaction_amount"`
	MasterStoreId     int    `validate:"required" json:"master_store_id"`
	Lang              string `validate:"required" json:"lang"`
	Extref            string `validate:"required" json:"extref"`
	InitiatorId       int    `validate:"required" json:"initiator_id"`
}

func (serv *CalculatorServiceImpl) GetPrice(ctx *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[Price](ctx.Body())
	if errData != nil {
		return model.Response{
			Status: fiber.StatusBadRequest,
			RC:     rc.FAILED.String(),
			Errors: errData,
		}
	}

	masterStore, err := serv.CalculatorRepository.PartnerMasterStoreId(req.MasterStoreId)
	if err != nil {
		return model.Response{
			Status: err.Status,
			RC:     err.RC,
			Errors: append([]model.ErrorData{}, *err),
		}
	}
	additionalData := AdditionalData{
		BankName:       masterStore.BankName,
		BankCode:       masterStore.BankCode,
		CustomerNumber: "123273784959",
		Remark:         "no",
	}

	initKey := os.Getenv("CRED_INIT_KEY")
	initName := os.Getenv("CRED_INIT_NAME")
	adminId := os.Getenv("CRED_ADMIN_ID")
	sendGetPrice := SendGetPrice{
		Signature:         utils.SignMD5("sendInq" + req.Extref + utils.SignSHA256(initKey) + initName),
		AdminId:           adminId,
		Extref:            "sendInq" + req.Extref,
		InitiatorId:       req.InitiatorId,
		PartnerId:         masterStore.PartnerId,
		CountryId:         masterStore.CountryId,
		TransactionAmount: req.TransactionAmount,
		AdditionalData:    &additionalData,
	}
	url := os.Getenv("BASE_URL")
	if req.Reverse != "yes" {
		url = url + "/sendmoneyinquiry"
	} else {
		url = url + "/sendmoneyinquiryreverse"
	}

	statusCode, mapBody, errResponse := utils.PostSendToUrl(sendGetPrice, url)
	if errResponse != nil {
		return model.Response{
			Status:  statusCode,
			RC:      rc.FAILED.String(),
			Message: "No Response from client",
			Errors:  append([]model.ErrorData{}, *errResponse),
		}
	}

	return model.Response{
		Status:  fiber.StatusOK,
		RC:      rc.SUCCESS.String(),
		Message: rc.SUCCESS.Message(),
		Extref:  req.Extref,
		Data:    mapBody,
	}
}

type AdditionalData struct {
	BankName       string `json:"bank_name"`
	BankCode       string `json:"bank_code"`
	CustomerNumber string `json:"customer_number"`
	Remark         string `json:"remark"`
}
type SendGetPrice struct {
	Signature         string          `json:"signature"`
	Extref            string          `json:"ext_ref"`
	AdminId           string          `json:"admin_id"`
	InitiatorId       int             `json:"initiator_id_sender"`
	PartnerId         string          `json:"partner_id_receiver"`
	CountryId         string          `json:"receiver_country"`
	TransactionAmount string          `json:"transaction_amount"`
	AdditionalData    *AdditionalData `json:"additional_data,omitempty"`
}
