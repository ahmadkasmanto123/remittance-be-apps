package service

import (
	"love-remittance-be-apps/core/rc"
	"love-remittance-be-apps/core/utils"
	"love-remittance-be-apps/lib/interfc"
	"love-remittance-be-apps/lib/model"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type InternationalServiceImpl struct {
	InternationalRepository interfc.InternationalRepository
}

func NewInternationalService(repository interfc.InternationalRepository) interfc.InternationalService {
	return &InternationalServiceImpl{
		InternationalRepository: repository,
	}
}

func (serv *InternationalServiceImpl) SourceOfPayment(c *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[string]](c.Body())
	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "error sini",
		}
	}

	datas, err := serv.InternationalRepository.GetSourceOfPayment()
	if err != nil {
		return model.Response{
			Status: err.Status,
			RC:     err.RC,
			Errors: append([]model.ErrorData{}, *err),
		}
	}
	var dynamicDatas []model.DynamicData[[]model.DataSOPayment]
	for _, val := range datas {
		indexGroup := -1
		// search existing type
		for i := 0; i < len(dynamicDatas); i++ {
			if dynamicDatas[i].DataType == strings.ToUpper(*val.SopType) { // break if type exist
				indexGroup = i // keep index group
				break
			}
		}
		// if type exist add to items
		if indexGroup != -1 {
			dynamicDatas[indexGroup].DataItems = append(dynamicDatas[indexGroup].DataItems, model.DataSOPayment{
				SopName:     val.SopName,
				SopId:       val.SopId,
				InitiatorId: val.InitiatorId,
			})
		} else { // if type not exist add new group with one item
			dynamicDatas = append(dynamicDatas, model.DynamicData[[]model.DataSOPayment]{
				DataType: strings.ToUpper(*val.SopType),
				DataItems: []model.DataSOPayment{{
					SopName:     val.SopName,
					SopId:       val.SopId,
					InitiatorId: val.InitiatorId,
				}},
			})
		}
	}
	return model.Response{
		Status:  fiber.StatusOK,
		RC:      rc.SUCCESS.String(),
		Message: rc.SUCCESS.Message(),
		Extref:  req.Extref,
		Data:    dynamicDatas,
	}
}

type sopData struct {
	SopId int `json:"sop_id"`
}

func (serv *InternationalServiceImpl) AvailableCountry(c *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[sopData]](c.Body())
	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "error sini",
		}
	}
	datas, err := serv.InternationalRepository.GetCountry(req.Request.SopId)
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

type countryData struct {
	CountryId int `json:"country_id"`
	SopId     int `json:"sop_id"`
}

func (serv *InternationalServiceImpl) DestinationByCountry(c *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[countryData]](c.Body())
	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "error sini",
		}
	}
	param := utils.GetParam(c)
	datas, err := serv.InternationalRepository.GetDestination(param, req.Request.CountryId, req.Request.SopId)
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

type getPrice struct {
	SopId         int    `validate:"required" json:"sop_id"`
	MasterStoreId int    `validate:"required" json:"master_store_id"`
	Reverse       string `validate:"required" json:"reverse"`
	Amount        string `validate:"required" json:"transaction_amount"`
}

func (serv *InternationalServiceImpl) GetPrice(c *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[getPrice]](c.Body())
	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "Please try again later",
		}
	}

	credent, err := serv.InternationalRepository.GetCredential(req.Request.MasterStoreId, req.Request.SopId)
	if err != nil {
		return model.Response{
			Status: err.Status,
			RC:     err.RC,
			Errors: append([]model.ErrorData{}, *err),
		}
	}

	additionalData := additionalData{
		BankName:       credent.BankName,
		BankCode:       credent.BankCode,
		CustomerNumber: "123273784959",
		Remark:         "no",
	}
	initKey := credent.InitiatorKey
	initName := credent.InitiatorUser
	adminId := credent.AdminId

	sendGetPrice := sendGetPriceInter{
		Signature:         utils.SignMD5("sendInq" + req.Extref + utils.SignSHA256(initKey) + initName),
		AdminId:           adminId,
		Extref:            "sendInq" + req.Extref,
		InitiatorId:       credent.InitiatorId,
		PartnerId:         credent.PartnerId,
		CountryId:         credent.CountryId,
		TransactionAmount: req.Request.Amount,
		AdditionalData:    &additionalData,
	}

	url := os.Getenv("BASE_URL")
	if req.Request.Reverse != "yes" {
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

	if mapBody["rc"].(string) != "00" {
		return model.Response{
			Status:  statusCode,
			RC:      rc.FAILED.String(),
			Message: "No Response from client",
		}
	}

	rateSender, _ := strconv.ParseFloat(mapBody["currency_rate_sender"].(string), 32)
	rateReceiver, _ := strconv.ParseFloat(mapBody["currency_rate_receiver"].(string), 32)
	rate := rateReceiver / rateSender

	resData := respDataGetPrice{
		AmountReceiveUsd:  mapBody["transaction_amount_dollar"].(string),
		AmountIdr:         mapBody["transaction_amount"].(string),
		AmountFeeIdr:      mapBody["transaction_fee"].(string),
		AmountTotalIdr:    mapBody["transaction_total_amount"].(string),
		CurrencyRate:      rate,
		AmountReceiveDest: mapBody["transaction_total_amount_receive"].(string),
	}
	return model.Response{
		Status:  fiber.StatusOK,
		RC:      rc.SUCCESS.String(),
		Message: rc.SUCCESS.Message(),
		Extref:  req.Extref,
		Data:    resData,
	}
}

type additionalData struct {
	BankName       string `json:"bank_name"`
	BankCode       string `json:"bank_code"`
	CustomerNumber string `json:"customer_number"`
	Remark         string `json:"remark"`
}
type sendGetPriceInter struct {
	Signature         string          `json:"signature"`
	Extref            string          `json:"ext_ref"`
	AdminId           string          `json:"admin_id"`
	InitiatorId       int             `json:"initiator_id_sender"`
	PartnerId         int             `json:"partner_id_receiver"`
	CountryId         int             `json:"receiver_country"`
	TransactionAmount string          `json:"transaction_amount"`
	AdditionalData    *additionalData `json:"additional_data,omitempty"`
}

type respDataGetPrice struct {
	AmountReceiveUsd  string  `json:"amount_receive_usd"`
	AmountIdr         string  `json:"transaction_amount"`
	AmountFeeIdr      string  `json:"transaction_fee"`
	AmountTotalIdr    string  `json:"transaction_total"`
	CurrencyRate      float64 `json:"rate"`
	AmountReceiveDest string  `json:"amount_receive_destination"`
}

func (serv *InternationalServiceImpl) CreatePaycode(c *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[model.Mandatory]](c.Body())
	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "Please try again later",
		}
	}

	credent, err := serv.InternationalRepository.GetCredential(req.Request.MasterStoreId, req.Request.SopId)
	if err != nil {
		return model.Response{
			Status: err.Status,
			RC:     err.RC,
			Errors: append([]model.ErrorData{}, *err),
		}
	}

	senderData, err := serv.InternationalRepository.DataAccountAll(req.Request.PhonePrefix, req.Request.Phone)
	if err != nil {
		return model.Response{
			Status: err.Status,
			RC:     err.RC,
			Errors: append([]model.ErrorData{}, *err),
		}
	}

	identitas, err := serv.InternationalRepository.GetIdentitas(req.Request.Receiver.IdentityType)
	if err != nil {
		return model.Response{
			Status: err.Status,
			RC:     err.RC,
			Errors: append([]model.ErrorData{}, *err),
		}
	}
	relations, err := serv.InternationalRepository.GetRelations(req.Request.Sender.RelationId)

	if err != nil {
		return model.Response{
			Status: err.Status,
			RC:     err.RC,
			Errors: append([]model.ErrorData{}, *err),
		}
	}
	purposes, err := serv.InternationalRepository.GetPurposes(req.Request.Sender.PurposeId)
	if err != nil {
		return model.Response{
			Status: err.Status,
			RC:     err.RC,
			Errors: append([]model.ErrorData{}, *err),
		}
	}
	fundings, err := serv.InternationalRepository.GetFunding(req.Request.Sender.SofundingId)
	if err != nil {
		return model.Response{
			Status: err.Status,
			RC:     err.RC,
			Errors: append([]model.ErrorData{}, *err),
		}
	}
	occupations, err := serv.InternationalRepository.GetOccupation(req.Request.Sender.OccupationId)
	if err != nil {
		return model.Response{
			Status: err.Status,
			RC:     err.RC,
			Errors: append([]model.ErrorData{}, *err),
		}
	}
	dataIdnKey := model.AllDataIdnKey{
		Identitas:   *identitas,
		Relations:   *relations,
		Purposes:    *purposes,
		Fundings:    *fundings,
		Occupations: *occupations,
	}
	sendReq, errss := PhlRequest(*req, *credent, *senderData, dataIdnKey)
	if errss != nil {
		return model.Response{
			Status: fiber.StatusBadRequest,
			RC:     rc.FAILED.String(),
			Errors: append([]model.ErrorData{}, model.ErrorData{
				Description: "reeeee",
			}),
		}
	}

	url := credent.UrlAdapter + "/create_paycode_sendmoney"
	statusCode, mapBody, errResponse := utils.PostSendToUrl(sendReq, url)
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

func PhlRequest(req model.DefaultRequest[model.Mandatory], credent model.Credential, senderData model.DataAccountDetail, data model.AllDataIdnKey) (model.ReqPhl, *model.ErrorData) {
	initKey := credent.InitiatorKey
	initName := credent.InitiatorUser
	adminId := credent.AdminId

	imgSelf, _ := utils.GetImgString(*senderData.ImgSelf)
	imgSign, _ := utils.GetImgString(*senderData.ImgSign)
	ImgIdentity, _ := utils.GetImgString(*senderData.ImgIdentity)

	var addtion model.AdditionalReqPhl

	if credent.CountryCode == "PHL" {
		addtion = model.AdditionalReqPhl{
			FirstName:            req.Request.Receiver.FirstName,
			LastName:             req.Request.Receiver.LastName,
			States:               req.Request.Receiver.States,
			PostalCodeReceiver:   req.Request.Receiver.PostalCode,
			IdentityTypeReceiver: data.Identitas.Key,
			PhoneReceiver:        req.Request.Receiver.Phone,
			CityReceiver:         req.Request.Receiver.CityName,
			AddressReceiver:      req.Request.Receiver.Address,
			NationalityReceiver:  credent.CountryCode,
			CustomerNumber:       req.Request.Additional.BenefiaciaryNumber,
			CustomerName:         req.Request.Additional.BenefiaciaryName,
			BankCode:             credent.BankCode,
			BankName:             credent.BankName,
			Remark:               req.Request.Additional.Remark,
			PostalCodeSender:     *senderData.PostalCode,
			OccupationSender:     data.Occupations.Key,
			SofundingSender:      data.Fundings.Key,
			PurposeSender:        data.Purposes.Key,
			RelationSender:       data.Relations.Key,
			CitySender:           *senderData.CityName,
			AmountReceiver:       req.Request.AmountReceive,
		}
	} else if credent.CountryCode == "AUS" {
		addtion = model.AdditionalReqPhl{
			FirstName:            req.Request.Receiver.FirstName,
			LastName:             req.Request.Receiver.LastName,
			States:               req.Request.Receiver.States,
			PostalCodeReceiver:   req.Request.Receiver.PostalCode,
			IdentityTypeReceiver: data.Identitas.Key,
			PhoneReceiver:        req.Request.Receiver.Phone,
			CityReceiver:         req.Request.Receiver.CityName,
			AddressReceiver:      req.Request.Receiver.Address,
			NationalityReceiver:  credent.CountryCode,
			CustomerNumber:       req.Request.Additional.BenefiaciaryNumber,
			CustomerName:         req.Request.Additional.BenefiaciaryName,
			BankCode:             credent.BankCode,
			BankName:             credent.BankName,
			BankBranchCode:       req.Request.Additional.BankBranchCode,
			Remark:               req.Request.Additional.Remark,
			PostalCodeSender:     *senderData.PostalCode,
			OccupationSender:     data.Occupations.Key,
			SofundingSender:      data.Fundings.Key,
			PurposeSender:        data.Purposes.Key,
			RelationSender:       data.Relations.Key,
			CitySender:           *senderData.CityName,
			AmountReceiver:       req.Request.AmountReceive,
		}
	} else if credent.CountryCode == "KHM" {
		addtion = model.AdditionalReqPhl{
			FirstName:            req.Request.Receiver.FirstName,
			LastName:             req.Request.Receiver.LastName,
			States:               req.Request.Receiver.States,
			PostalCodeReceiver:   req.Request.Receiver.PostalCode,
			IdentityTypeReceiver: data.Identitas.Key,
			PhoneReceiver:        req.Request.Receiver.Phone,
			CityReceiver:         req.Request.Receiver.CityName,
			AddressReceiver:      req.Request.Receiver.Address,
			NationalityReceiver:  credent.CountryCode,
			CustomerNumber:       req.Request.Additional.BenefiaciaryNumber,
			CustomerName:         req.Request.Additional.BenefiaciaryName,
			BankCode:             credent.BankCode,
			BankName:             credent.BankName,
			BankBranchCode:       req.Request.Additional.BankBranchCode,
			Remark:               req.Request.Additional.Remark,
			PostalCodeSender:     *senderData.PostalCode,
			OccupationSender:     data.Occupations.Key,
			SofundingSender:      data.Fundings.Key,
			PurposeSender:        data.Purposes.Key,
			RelationSender:       data.Relations.Key,
			CitySender:           *senderData.CityName,
			AmountReceiver:       req.Request.AmountReceive,
		}
	}

	sendReq := model.ReqPhl{
		Signature:   utils.SignMD5("sendInq" + req.Extref + utils.SignSHA256(initKey) + initName),
		Extref:      req.Extref,
		Amount:      req.Request.Amount,
		Phone:       req.Request.PhonePrefix + req.Request.Phone,
		PartnerId:   credent.PartnerId,
		Data1:       req.Request.Additional.Remark,
		Data2:       "",
		ImgSelf:     imgSelf,
		ImgSign:     imgSign,
		ImgIdentity: ImgIdentity,
		Sender: model.SenderReqPhl{
			Name:           senderData.AccountName,
			Phone:          req.Request.PhonePrefix + req.Request.Phone,
			Address:        *senderData.Address,
			OccupationId:   data.Occupations.Id,
			SopId:          req.Request.SopId,
			Gender:         *senderData.Gender,
			PurposeId:      data.Purposes.Id,
			SofundingId:    data.Fundings.Id,
			DOB:            senderData.DOB.Format("2006-01-02"),
			POB:            *senderData.POB,
			IdentityTypeId: *senderData.IdentityTypeId,
			IdentityNumber: *senderData.IdentityNumber,
			IdentityExp:    "2500-01-02",
			AdminId:        adminId,
			City:           *senderData.CityId,
			Data1:          req.Request.Additional.Remark,
		},
		Receiver: model.ReceiverReqPhl{
			FirstName:      req.Request.Receiver.FirstName,
			LastName:       req.Request.Receiver.LastName,
			Address:        req.Request.Receiver.Address,
			City:           credent.CityReceiver,
			Name:           req.Request.Receiver.FirstName + " " + req.Request.Receiver.LastName,
			Phone:          req.Request.Receiver.Phone,
			IdentityNumber: req.Request.Receiver.IdentityNumber,
			IdentityTypeID: data.Identitas.Id,
		},
		Additional: addtion,
	}
	return sendReq, nil
}
