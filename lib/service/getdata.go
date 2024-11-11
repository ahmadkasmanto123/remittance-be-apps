package service

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"love-remittance-be-apps/core/rc"
	"love-remittance-be-apps/core/utils"
	"love-remittance-be-apps/lib/interfc"
	"love-remittance-be-apps/lib/model"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type GetDataAllServiceImpl struct {
	GetDataAllRepository interfc.GetDataAllRepository
}

func NewGetDataAllService(repos interfc.GetDataAllRepository) interfc.GetDataAllService {
	return &GetDataAllServiceImpl{
		GetDataAllRepository: repos,
	}
}

func (serv *GetDataAllServiceImpl) GetImage(ctx *fiber.Ctx) error {
	imgName := ctx.Query("img_name")
	conn, errCon := utils.FtpConnection()
	if errCon != nil {
		log.Println(errCon.Description)
	}
	desFile := "." + os.Getenv("FTP_PATH") + "/" + imgName
	files, err := conn.Retr(desFile)
	if err != nil {
		log.Println(err.Error())
	}

	imgByte, err := io.ReadAll(files)
	defer files.Close()
	if err != nil {
		log.Println(err.Error())
	}
	contentType := http.DetectContentType(imgByte)

	ctx.Attachment("./file/" + imgName)
	ctx.Set("Content-Type", contentType)
	ctx.Status(fiber.StatusOK)

	return ctx.SendStream(bytes.NewReader(imgByte))
}
func (serv *GetDataAllServiceImpl) GetCountry(c *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[string]](c.Body())
	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "error sini",
		}
	}
	param := utils.GetParam(c)
	datas, err := serv.GetDataAllRepository.GetCountry(param)
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
		Data:    datas,
	}
}
func (serv *GetDataAllServiceImpl) GetProvince(c *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[string]](c.Body())
	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "error sini",
		}
	}
	param := utils.GetParam(c)
	datas, pagination, error := serv.GetDataAllRepository.GetProvince(param)
	if error != nil {
		return model.Response{
			Status: error.Status,
			RC:     error.RC,
			Errors: append([]model.ErrorData{}, *error),
		}
	}

	return model.Response{
		Status:     fiber.StatusOK,
		RC:         rc.SUCCESS.String(),
		Message:    rc.SUCCESS.Message(),
		Extref:     req.Extref,
		Data:       datas,
		Pagination: pagination,
	}
}
func (serv *GetDataAllServiceImpl) GetCity(c *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[model.Province]](c.Body())
	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "error sini",
		}
	}
	param := utils.GetParam(c)
	datas, pagination, error := serv.GetDataAllRepository.GetCity(param, req.Request.ProvinceId)
	if error != nil {
		return model.Response{
			Status: error.Status,
			RC:     error.RC,
			Errors: append([]model.ErrorData{}, *error),
		}
	}

	return model.Response{
		Status:     fiber.StatusOK,
		RC:         rc.SUCCESS.String(),
		Message:    rc.SUCCESS.Message(),
		Extref:     req.Extref,
		Data:       datas,
		Pagination: pagination,
	}
}
func (serv *GetDataAllServiceImpl) GetOccupation(c *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[string]](c.Body())
	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "error sini",
		}
	}
	datas, error := serv.GetDataAllRepository.GetOccupation()
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
		Data:    datas}
}
func (serv *GetDataAllServiceImpl) GetIdentityType(c *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[string]](c.Body())
	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "error sini",
		}
	}
	datas, error := serv.GetDataAllRepository.GetIdentityType()
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
		Data:    datas}
}

type accountData struct {
	Phone       string `json:"phone" validate:"required"`
	PhonePrefix string `json:"phone_prefix" validate:"required"`
}

type accountDetail struct {
	AccountId         int        `json:"id"`
	AccountStatutId   int        `json:"status_id"`
	AccountStatusName string     `json:"status_name"`
	IdentityTypeId    *int       `json:"identity_type_id,omitempty"`
	IdentityTypeName  *string    `json:"identity_type_name,omitempty"`
	IdentityNumber    *string    `json:"identity_number,omitempty"`
	FirstName         string     `json:"first_name"`
	LastName          string     `json:"last_name"`
	AdminEmail        string     `json:"admin_email"`
	DeviceId          string     `json:"device_id"`
	CityId            *int       `json:"city_id,omitempty"`
	CityName          *string    `json:"city_name,omitempty"`
	OccupationId      *int       `json:"occupation_id,omitempty"`
	OccupationName    *string    `json:"occupation_name,omitempty"`
	Address           *string    `json:"address,omitempty"`
	POB               *string    `json:"pob,omitempty"`
	DOB               *time.Time `json:"dob,omitempty"`
	Gender            *string    `json:"gender,omitempty"`
	PostalCode        *string    `json:"postal_code,omitempty"`
	ImgSelf           *string    `json:"img_self,omitempty"`
}

func (serv *GetDataAllServiceImpl) GetAccountDetail(c *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[accountData]](c.Body())
	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "error sini",
		}
	}
	datas, error := serv.GetDataAllRepository.DataAccountAll(req.Request.PhonePrefix, req.Request.Phone)
	if error != nil {
		return model.Response{
			Status: error.Status,
			RC:     error.RC,
			Errors: append([]model.ErrorData{}, *error),
		}
	}

	accName := strings.Split(datas.AccountName, " ")
	var firstName string
	var lastName string
	if len(accName) > 1 {
		firstName = accName[0]
		lastName = accName[1]
	}

	result := accountDetail{
		AccountId:         datas.AccountId,
		AccountStatutId:   datas.AccountStatutId,
		AccountStatusName: datas.AccountStatusName,
		IdentityTypeId:    datas.IdentityTypeId,
		IdentityTypeName:  datas.IdentityTypeName,
		IdentityNumber:    datas.IdentityNumber,
		FirstName:         firstName,
		LastName:          lastName,
		AdminEmail:        datas.AdminEmail,
		DeviceId:          datas.DeviceId,
		CityId:            datas.CityId,
		CityName:          datas.CityName,
		OccupationId:      datas.OccupationId,
		OccupationName:    datas.OccupationName,
		Address:           datas.Address,
		POB:               datas.POB,
		DOB:               datas.DOB,
		Gender:            datas.Gender,
		PostalCode:        datas.PostalCode,
		ImgSelf:           datas.ImgSelf,
	}

	return model.Response{
		Status:  fiber.StatusOK,
		RC:      rc.SUCCESS.String(),
		Message: rc.SUCCESS.Message(),
		Extref:  req.Extref,
		Data:    result}
}

func (serv *GetDataAllServiceImpl) GetSoFund(c *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[string]](c.Body())
	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "error sini",
		}
	}
	datas, error := serv.GetDataAllRepository.GetSoFund()
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
		Data:    datas}
}

func (serv *GetDataAllServiceImpl) GetPurpose(c *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[string]](c.Body())
	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "error sini",
		}
	}
	datas, error := serv.GetDataAllRepository.GetPurpose()
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
		Data:    datas}
}

func (serv *GetDataAllServiceImpl) GetRelations(c *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[string]](c.Body())
	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "error sini",
		}
	}
	datas, error := serv.GetDataAllRepository.GetRelations()
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
		Data:    datas}
}

type country struct {
	CountryCode string `json:"country_code"`
}

func (serv *GetDataAllServiceImpl) GetFromSender(c *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[country]](c.Body())
	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "error sini",
		}
	}
	// Let's first read the `config.json` file
	nameFile := utils.GetFileConfig(req.Request.CountryCode, req.Lang)
	content, err := os.ReadFile(nameFile)
	if err != nil {
		return model.Response{
			Status: fiber.StatusBadRequest,
			RC:     rc.FAILED.String(),
			Errors: append([]model.ErrorData{}, model.ErrorData{
				Description: err.Error(),
			}),
			Message: "Form tidak tersedia",
		}
	}

	// Now let's unmarshall the data into `payload`
	var payload map[string]interface{}
	err = json.Unmarshal(content, &payload)
	if err != nil {
		return model.Response{
			Status: fiber.StatusBadRequest,
			RC:     rc.FAILED.String(),
			Errors: append([]model.ErrorData{}, model.ErrorData{
				Description: err.Error(),
			}),
			Message: "Error during Unmarshal()",
		}
	}

	data := utils.StringToMap(payload["data"])
	return model.Response{
		Status:  fiber.StatusOK,
		RC:      rc.SUCCESS.String(),
		Message: rc.SUCCESS.Message(),
		Extref:  req.Extref,
		Data:    data["sender"],
	}
}

func (serv *GetDataAllServiceImpl) GetFromReceiver(c *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[country]](c.Body())
	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "error sini",
		}
	}
	// Let's first read the `config.json` file
	nameFile := utils.GetFileConfig(req.Request.CountryCode, req.Lang)
	content, err := os.ReadFile(nameFile)
	if err != nil {
		return model.Response{
			Status: fiber.StatusBadRequest,
			RC:     rc.FAILED.String(),
			Errors: append([]model.ErrorData{}, model.ErrorData{
				Description: err.Error(),
			}),
			Message: "Form tidak tersedia",
		}
	}

	// Now let's unmarshall the data into `payload`
	var payload map[string]interface{}
	err = json.Unmarshal(content, &payload)
	if err != nil {
		return model.Response{
			Status: fiber.StatusBadRequest,
			RC:     rc.FAILED.String(),
			Errors: append([]model.ErrorData{}, model.ErrorData{
				Description: err.Error(),
			}),
			Message: "Error during Unmarshal()",
		}
	}

	data := utils.StringToMap(payload["data"])
	return model.Response{
		Status:  fiber.StatusOK,
		RC:      rc.SUCCESS.String(),
		Message: rc.SUCCESS.Message(),
		Extref:  req.Extref,
		Data:    data["receiver"],
	}
}

func (serv *GetDataAllServiceImpl) GetFromAdditional(c *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[country]](c.Body())
	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "error sini",
		}
	}
	// Let's first read the `config.json` file
	nameFile := utils.GetFileConfig(req.Request.CountryCode, req.Lang)
	content, err := os.ReadFile(nameFile)
	if err != nil {
		return model.Response{
			Status: fiber.StatusBadRequest,
			RC:     rc.FAILED.String(),
			Errors: append([]model.ErrorData{}, model.ErrorData{
				Description: err.Error(),
			}),
			Message: "Form tidak tersedia",
		}
	}

	// Now let's unmarshall the data into `payload`
	var payload map[string]interface{}
	err = json.Unmarshal(content, &payload)
	if err != nil {
		return model.Response{
			Status: fiber.StatusBadRequest,
			RC:     rc.FAILED.String(),
			Errors: append([]model.ErrorData{}, model.ErrorData{
				Description: err.Error(),
			}),
			Message: "Error during Unmarshal()",
		}
	}

	data := utils.StringToMap(payload["data"])
	return model.Response{
		Status:  fiber.StatusOK,
		RC:      rc.SUCCESS.String(),
		Message: rc.SUCCESS.Message(),
		Extref:  req.Extref,
		Data:    data["additional"],
	}
}

type master struct {
	MasterStoreId int `json:"master_store_id"`
}

func (serv *GetDataAllServiceImpl) GetOccupationIntl(c *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[master]](c.Body())
	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "error sini",
		}
	}
	datas, error := serv.GetDataAllRepository.GetOccupationIntl(req.Request.MasterStoreId)
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
		Data:    datas}
}
func (serv *GetDataAllServiceImpl) GetIdentityTypeIntl(c *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[master]](c.Body())
	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "error sini",
		}
	}
	datas, error := serv.GetDataAllRepository.GetIdentityTypeIntl(req.Request.MasterStoreId)
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
		Data:    datas}
}

func (serv *GetDataAllServiceImpl) GetSoFundIntl(c *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[master]](c.Body())
	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "error sini",
		}
	}
	datas, error := serv.GetDataAllRepository.GetSoFundIntl(req.Request.MasterStoreId)
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
		Data:    datas}
}

func (serv *GetDataAllServiceImpl) GetPurposeIntl(c *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[master]](c.Body())
	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "error sini",
		}
	}
	datas, error := serv.GetDataAllRepository.GetPurposeIntl(req.Request.MasterStoreId)
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
		Data:    datas}
}

func (serv *GetDataAllServiceImpl) GetRelationsIntl(c *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[master]](c.Body())
	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "error sini",
		}
	}
	datas, error := serv.GetDataAllRepository.GetRelationsIntl(req.Request.MasterStoreId)
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
		Data:    datas}
}

func (serv *GetDataAllServiceImpl) GetHistory(c *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[model.DataPhone]](c.Body())
	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "error sini",
		}
	}
	datas, error := serv.GetDataAllRepository.GetHistory(req.Request.PhonePrefix + req.Request.Phone)
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
		Data:    datas,
	}
}

type onlyTransaction struct {
	TransactionId int `json:"transaction_id"`
}

func (serv *GetDataAllServiceImpl) GetHistoryDetail(c *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[onlyTransaction]](c.Body())
	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "error sini",
		}
	}
	datas, error := serv.GetDataAllRepository.GetHistoryDetail(req.Request.TransactionId)
	if error != nil {
		return model.Response{
			Status: error.Status,
			RC:     error.RC,
			Errors: append([]model.ErrorData{}, *error),
		}
	}

	url := datas.UrlAdapter + "/check_voucher"
	sendReq := sendCheckVoucher{
		TransactionId: req.Request.TransactionId,
		ExtRef:        req.Extref,
	}
	statusCode, mapBody, errResponse := utils.PostSendToUrl(sendReq, url)
	if errResponse != nil {
		return model.Response{
			Status:  statusCode,
			RC:      rc.FAILED.String(),
			Message: "No Response from client",
			Errors:  append([]model.ErrorData{}, *errResponse),
		}
	}

	payload, err := utils.GetParameterNotes("NOTE", req.Lang)
	if err != nil {
		return model.Response{
			Status:  statusCode,
			RC:      rc.FAILED.String(),
			Message: "Error Get Foot Note",
			Errors:  append([]model.ErrorData{}, *err),
		}
	}

	var dataTrx map[string]interface{}
	var addData map[string]interface{}
	amountReceive := ""
	bankBranchCode := ""
	if mapBody["data_transaction"] != nil {
		dataTrx = utils.StringToMap(mapBody["data_transaction"])
		if dataTrx["additional_data"] != nil {
			addData = utils.StringToMap(dataTrx["additional_data"])
			if addData["amount_receiver"] != nil {
				amountReceive = addData["amount_receiver"].(string)
			}
			if addData["bank_branch_code"] != nil {
				bankBranchCode = addData["bank_branch_code"].(string)
			}
		}
	}

	var resp any

	if datas.TypeId == 1 {
		resp = model.VoucherActive{
			BankName:           addData["bank_name"].(string),
			BankBranchCode:     bankBranchCode,
			BenefiaciaryName:   addData["customer_name"].(string),
			BenefiaciaryNumber: addData["customer_number"].(string),
			Remark:             addData["remark"].(string),
			Amount:             mapBody["amount"].(string),
			AmountFee:          mapBody["admin"].(string),
			AmountTotal:        mapBody["total_amount"].(string),
			AmountReceive:      amountReceive,
			PaycodeNumber:      mapBody["payment_code"].(string),
			Journey:            "",
		}
	} else if datas.TypeId == 2 || datas.TypeId == 4 {
		if datas.TypeTrx == "BANK" {
			resp = model.HistoryDetail{
				BankName:           addData["bank_name"].(string),
				BankBranchCode:     bankBranchCode,
				BenefiaciaryName:   addData["customer_name"].(string),
				BenefiaciaryNumber: addData["customer_number"].(string),
				Remark:             addData["remark"].(string),
				Amount:             mapBody["amount"].(string),
				AmountFee:          mapBody["admin"].(string),
				AmountTotal:        mapBody["total_amount"].(string),
				AmountReceive:      amountReceive,
				Rate:               "",
				SenderName:         dataTrx["sender_name"].(string),
				SenderAddress:      dataTrx["sender_address"].(string),
				TransactionId:      req.Request.TransactionId,
				DateTrx:            datas.TimeTrx,
				TypeTrx:            datas.TypeTrx,
				ReceiverName:       datas.ReceiverName,
				ReceiverPhone:      datas.ReceiverPhone,
				ReceiverCountry:    datas.ReceiverCountry,
				FootNote:           payload["foot_note"].(string),
			}
		} else {
			resp = model.HistoryDetail{
				BankName:           addData["bank_name"].(string),
				BankBranchCode:     bankBranchCode,
				BenefiaciaryName:   addData["customer_name"].(string),
				BenefiaciaryNumber: addData["customer_number"].(string),
				Remark:             addData["remark"].(string),
				Amount:             mapBody["amount"].(string),
				AmountFee:          mapBody["admin"].(string),
				AmountTotal:        mapBody["total_amount"].(string),
				AmountReceive:      amountReceive,
				Rate:               "",
				SenderName:         dataTrx["sender_name"].(string),
				SenderAddress:      dataTrx["sender_address"].(string),
				TransactionId:      req.Request.TransactionId,
				DateTrx:            datas.TimeTrx,
				TypeTrx:            datas.TypeTrx,
				ReceiverName:       datas.ReceiverName,
				ReceiverPhone:      datas.ReceiverPhone,
				ReceiverCountry:    datas.ReceiverCountry,
				FootNote:           payload["foot_note"].(string),
			}
		}
	} else {
		resp = model.HistoryDetail{
			BankName:           addData["bank_name"].(string),
			BankBranchCode:     bankBranchCode,
			BenefiaciaryName:   addData["customer_name"].(string),
			BenefiaciaryNumber: addData["customer_number"].(string),
			Remark:             addData["remark"].(string),
			Amount:             mapBody["amount"].(string),
			AmountFee:          mapBody["admin"].(string),
			AmountTotal:        mapBody["total_amount"].(string),
			AmountReceive:      amountReceive,
			Rate:               "",
			SenderName:         dataTrx["sender_name"].(string),
			SenderAddress:      dataTrx["sender_address"].(string),
			TransactionId:      req.Request.TransactionId,
			DateTrx:            datas.TimeTrx,
			TypeTrx:            datas.TypeTrx,
			ReceiverName:       datas.ReceiverName,
			ReceiverPhone:      datas.ReceiverPhone,
			ReceiverCountry:    datas.ReceiverCountry,
			FootNote:           payload["foot_note"].(string),
		}
	}

	return model.Response{
		Status:  fiber.StatusOK,
		RC:      rc.SUCCESS.String(),
		Message: rc.SUCCESS.Message(),
		Extref:  req.Extref,
		Data:    resp,
	}
}

type sendCheckVoucher struct {
	TransactionId int    `json:"transaction_id"`
	ExtRef        string `json:"ext_ref"`
}
