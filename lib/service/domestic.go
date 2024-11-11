package service

import (
	"love-remittance-be-apps/core/rc"
	"love-remittance-be-apps/core/utils"
	"love-remittance-be-apps/lib/interfc"
	"love-remittance-be-apps/lib/model"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type DomesticServiceImpl struct {
	DomesticRepository interfc.DomesticRepository
}

func NewDomesticService(repository interfc.DomesticRepository) interfc.DomesticService {
	return &DomesticServiceImpl{
		DomesticRepository: repository,
	}
}

// type requestStandar struct {
// 	Extref string `validate:"required" json:"extref"`
// 	Lang   string `validate:"required" json:"lang"`
// }

func (serv *DomesticServiceImpl) Sopayment(c *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[string]](c.Body())
	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "error sini",
		}
	}
	datas, error := serv.DomesticRepository.GetSourceOfPayment()
	if error != nil {
		return model.Response{
			Status: error.Status,
			RC:     error.RC,
			Errors: append([]model.ErrorData{}, *error),
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

type destination struct {
	SopId int `validate:"required" json:"sop_id"`
}

func (serv *DomesticServiceImpl) Destination(c *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[destination]](c.Body())
	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "error sini",
		}
	}
	param := utils.GetParam(c)
	datas, pagination, error := serv.DomesticRepository.GetDestination(param, req.Request.SopId)
	if error != nil {
		return model.Response{
			Status: error.Status,
			RC:     error.RC,
			Errors: append([]model.ErrorData{}, *error),
		}
	}

	var dynamicDatas []model.DynamicData[[]model.DataDestination]
	for _, val := range datas {
		indexGroup := -1
		// search existing type
		for i := 0; i < len(dynamicDatas); i++ {
			if dynamicDatas[i].DataType == strings.ToUpper(val.Types) { // break if type exist
				indexGroup = i // keep index group
				break
			}
		}
		// if type exist add to items
		if indexGroup != -1 {
			dynamicDatas[indexGroup].DataItems = append(dynamicDatas[indexGroup].DataItems, model.DataDestination{
				BankName:      val.BankName,
				BankCode:      val.BankCode,
				MasterStoreId: val.MasterStoreId,
			})
		} else { // if type not exist add new group with one item
			dynamicDatas = append(dynamicDatas, model.DynamicData[[]model.DataDestination]{
				DataType: strings.ToUpper(val.Types),
				DataItems: []model.DataDestination{{
					BankName:      val.BankName,
					BankCode:      val.BankCode,
					MasterStoreId: val.MasterStoreId,
				}},
			})
		}
	}

	return model.Response{
		Status:     fiber.StatusOK,
		RC:         rc.SUCCESS.String(),
		Message:    rc.SUCCESS.Message(),
		Extref:     req.Extref,
		Data:       dynamicDatas,
		Pagination: pagination,
	}
}

type price struct {
	TransactionAmount string `validate:"required" json:"transaction_amount"`
	MasterStoreId     int    `validate:"required" json:"master_store_id"`
	SopId             int    `validate:"required" json:"sop_id"`
}

func (serv *DomesticServiceImpl) GetPrice(ctx *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[price]](ctx.Body())
	if errData != nil {
		return model.Response{
			Status: fiber.StatusBadRequest,
			RC:     rc.FAILED.String(),
			Errors: errData,
		}
	}
	masterStore, err := serv.DomesticRepository.PartnerMasterStoreId(req.Request.MasterStoreId)
	if err != nil {
		return model.Response{
			Status: err.Status,
			RC:     err.RC,
			Errors: append([]model.ErrorData{}, *err),
		}
	}
	sopData, errss := serv.DomesticRepository.DataInitiator(req.Request.SopId)
	if errss != nil {
		return model.Response{
			Status: errss.Status,
			RC:     errss.RC,
			Errors: append([]model.ErrorData{}, *errss),
		}
	}

	initKey := os.Getenv("CRED_INIT_KEY")
	initName := os.Getenv("CRED_INIT_NAME")
	adminId := os.Getenv("CRED_ADMIN_ID")

	sendGetPrice := sendGetPrice{
		Signature:         utils.SignMD5("sendInq" + req.Extref + utils.SignSHA256(initKey) + initName),
		AdminId:           adminId,
		Extref:            "sendInq" + req.Extref,
		InitiatorId:       *sopData.InitiatorId,
		PartnerId:         masterStore.PartnerId,
		CountryId:         masterStore.CountryId,
		TransactionAmount: req.Request.TransactionAmount,
	}
	url := os.Getenv("BASE_URL") + "/sendmoneyinquiry"
	statusCode, mapBody, errResponse := utils.PostSendToUrl(sendGetPrice, url)
	if errResponse != nil {
		return model.Response{
			Status:  statusCode,
			RC:      rc.FAILED.String(),
			Message: "No Response from client",
			Errors:  append([]model.ErrorData{}, *errResponse),
		}
	}

	resGetPrice := resGetPrice{
		TransactionAmount: mapBody["transaction_amount"].(string),
		TransactionFee:    mapBody["transaction_fee"].(string),
		TransactionTotal:  mapBody["transaction_total_amount"].(string),
		CurrencyCode:      mapBody["currency_code_receiver"].(string),
		CurrencyName:      mapBody["currency_name_receiver"].(string),
	}
	return model.Response{
		Status:  fiber.StatusOK,
		RC:      rc.SUCCESS.String(),
		Message: rc.SUCCESS.Message(),
		Extref:  req.Extref,
		Data:    resGetPrice,
	}
}

type sendGetPrice struct {
	Signature         string `json:"signature"`
	Extref            string `json:"ext_ref"`
	AdminId           string `json:"admin_id"`
	InitiatorId       int    `json:"initiator_id_sender"`
	PartnerId         string `json:"partner_id_receiver"`
	CountryId         string `json:"receiver_country"`
	TransactionAmount string `json:"transaction_amount"`
}
type resGetPrice struct {
	TransactionAmount string `json:"transaction_amount"`
	TransactionFee    string `json:"transaction_fee"`
	TransactionTotal  string `json:"transaction_total"`
	CurrencyCode      string `json:"currency_code"`
	CurrencyName      string `json:"currency_name"`
}
type request struct {
	SopId         int           `json:"sop_id"`
	MasterStoreId int           `json:"master_store_id"`
	Amount        string        `json:"transaction_amount"`
	Phone         string        `json:"phone"`
	PhonePrefix   string        `json:"phone_prefix"`
	Sender        reqSender     `json:"sender"`
	Receiver      reqReceiver   `json:"receiver"`
	Additional    reqAdditional `json:"additional"`
}
type reqSender struct {
	SofId        int `json:"sofunding_id"`
	PurposeId    int `json:"purpose_id"`
	OccupationId int `json:"occupation_id"`
}
type reqReceiver struct {
	CityId      int    `json:"city_id"`
	Phone       string `json:"phone"`
	PhonePrefix string `json:"phone_prefix"`
}
type reqAdditional struct {
	Remark      string `json:"remark"`
	BenerNumber string `json:"beneficiary_number"`
	BenerName   string `json:"beneficiary_name"`
}

func (serv *DomesticServiceImpl) CheckAccount(ctx *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[request]](ctx.Body())
	if errData != nil {
		return model.Response{
			Status: fiber.StatusBadRequest,
			RC:     rc.FAILED.String(),
			Errors: errData,
		}
	}
	masterStore, err := serv.DomesticRepository.PartnerMasterStoreId(req.Request.MasterStoreId)
	if err != nil {
		return model.Response{
			Status: err.Status,
			RC:     err.RC,
			Errors: append([]model.ErrorData{}, *err),
		}
	}
	senderData, errss := serv.DomesticRepository.DataAccountAll(req.Request.PhonePrefix, req.Request.Phone)
	if errss != nil {
		return model.Response{
			Status: errss.Status,
			RC:     errss.RC,
			Errors: append([]model.ErrorData{}, *errss),
		}
	}

	initKey := os.Getenv("CRED_INIT_KEY")
	initName := os.Getenv("CRED_INIT_NAME")
	adminId := os.Getenv("CRED_ADMIN_ID")

	sendCheckAccount := sendCheckAccount{
		Signature: utils.SignMD5(req.Extref + utils.SignSHA256(initKey) + initName),
		PartnerId: masterStore.PartnerId,
		AdminId:   adminId,
		FullUri:   "inquiry_sendmoney",
		Extref:    req.Extref,
		CurlData: curlData{
			Amount: req.Request.Amount,
			Sender: sendSender{
				Name:         senderData.AccountName,
				Phone:        req.Request.PhonePrefix + req.Request.Phone,
				Address:      *senderData.Address,
				OccupationId: *senderData.OccupationId,
				SopId:        req.Request.SopId,
				Gender:       *senderData.Gender,
				CityId:       *senderData.CityId,
				PurposeId:    req.Request.Sender.PurposeId,
				SofundingId:  req.Request.Sender.SofId,
				DOB:          senderData.DOB.Format("2006-01-02"),
				POB:          *senderData.POB,
				IdTypeId:     *senderData.IdentityTypeId,
				IdNumber:     *senderData.IdentityNumber,
				ProvinceName: "",
			},
			Receiver: sendReceiver{
				Phone:          req.Request.Receiver.PhonePrefix + req.Request.Receiver.Phone,
				CityId:         req.Request.Receiver.CityId,
				Address:        "",
				BankCode:       masterStore.BankCode,
				BankName:       masterStore.BankName,
				IdTypeID:       "",
				IdNumber:       "",
				CustomerNumber: req.Request.Additional.BenerNumber,
				ProvinceName:   "",
			},
			Additional: sendAdditional{
				CustomerNumber:   req.Request.Additional.BenerNumber,
				BankCode:         masterStore.BankCode,
				BankName:         masterStore.BankName,
				OccupationSender: "employee",
				Remark:           req.Request.Additional.Remark,
			},
		},
	}
	url := os.Getenv("BASE_URL") + "/freeUrl"
	// fmt.Printf("%v", sendCheckAccount)
	statusCode, mapBody, errResponse := utils.PostSendToUrl(sendCheckAccount, url)
	if errResponse != nil {
		return model.Response{
			Status:  statusCode,
			RC:      rc.FAILED.String(),
			Message: "No Response from client",
			Errors:  append([]model.ErrorData{}, *errResponse),
		}
	}
	resultResp := model.Response{
		Status:  fiber.StatusOK,
		RC:      rc.FAILED.String(),
		Message: rc.FAILED.Message(),
		Extref:  req.Extref,
		Data: respData{
			CustomerNumber: req.Request.Additional.BenerNumber,
			CustomerName:   "",
			BankCode:       masterStore.BankCode,
			BankName:       masterStore.BankName,
		},
	}
	// check if data not exist
	if mapBody["rc"].(string) != "00" {
		if mapBody["rc"].(string) == "2039" {
			if mapBody["jObjGwResp"] != nil {
				resp := utils.StringToMap(mapBody["jObjGwResp"])
				if resp["rc"].(string) == "3005" {
					return model.Response{
						Status:  fiber.StatusOK,
						RC:      rc.FAILED.String(),
						Message: resp["message"].(string),
						Extref:  req.Extref,
						Data: respData{
							CustomerNumber: req.Request.Additional.BenerNumber,
							CustomerName:   "PENDING",
							BankCode:       masterStore.BankCode,
							BankName:       masterStore.BankName,
						},
					}
				} else if resp["rc"].(string) == "3011" {
					return model.Response{
						Status:  fiber.StatusOK,
						RC:      rc.FAILED.String(),
						Message: "Invalid account number",
						Extref:  req.Extref,
						Data: respData{
							CustomerNumber: req.Request.Additional.BenerNumber,
							CustomerName:   "",
							BankCode:       masterStore.BankCode,
							BankName:       masterStore.BankName,
						},
					}
				} else {
					return resultResp
				}
			} else {
				return resultResp
			}
		} else {
			return resultResp
		}
	}

	resp := utils.StringToMap(mapBody["jObjGwResp"])
	return model.Response{
		Status:  fiber.StatusOK,
		RC:      rc.SUCCESS.String(),
		Message: rc.SUCCESS.Message(),
		Extref:  req.Extref,
		Data: respData{
			CustomerNumber: resp["account_number"].(string),
			CustomerName:   resp["account_holder"].(string),
			BankName:       resp["bank_name"].(string),
			BankCode:       resp["bank_code"].(string),
		},
	}
}

type sendCheckAccount struct {
	Signature string   `json:"signature"`
	Extref    string   `json:"ext_ref"`
	AdminId   string   `json:"admin_id"`
	PartnerId string   `json:"partner_id_receiver"`
	FullUri   string   `json:"full_uri"`
	CurlData  curlData `json:"curl_data"`
}
type curlData struct {
	Amount     string         `json:"amount"`
	Sender     sendSender     `json:"sender"`
	Receiver   sendReceiver   `json:"receiver"`
	Additional sendAdditional `json:"additional_data"`
}
type sendSender struct {
	Name           string `json:"name"`
	Phone          string `json:"phone_number"`
	Address        string `json:"address"`
	OccupationId   int    `json:"occupation"`
	SopId          int    `json:"source_of_payment"`
	Gender         string `json:"gender"`
	CityId         int    `json:"city_name,omitempty"`
	PurposeId      int    `json:"purpose"`
	SofundingId    int    `json:"source_of_fund"`
	DOB            string `json:"dob"`
	POB            string `json:"pob"`
	IdTypeId       int    `json:"id_type,omitempty"`
	IdNumber       string `json:"id_number,omitempty"`
	IdExp          string `json:"id_exp,omitempty"`
	ProvinceName   string `json:"province_name,omitempty"`
	IdentityTypeId int    `json:"identity_type"`
	IdentityNumber string `json:"identity_number"`
	IdentityExp    string `json:"identity_exp"`
	AdminId        int    `json:"admin_id"`
	City           int    `json:"city"`
}
type sendReceiver struct {
	Address        string  `json:"address"`
	CityId         int     `json:"city_name,omitempty"`
	Phone          string  `json:"phone_number"`
	CustomerNumber string  `json:"account_number"`
	BankCode       string  `json:"bank_code"`
	BankName       string  `json:"bank_name"`
	Remark         string  `json:"remark"`
	CustomerName   *string `json:"name,omitempty"`
	IdNumber       string  `json:"id_number"`
	IdTypeID       string  `json:"id_type"`
	IdentityNumber string  `json:"identity_number"`
	IdentityTypeID string  `json:"identity_type"`
	ProvinceName   string  `json:"province_name"`
	City           int     `json:"city"`
}
type sendAdditional struct {
	CustomerNumber   string `json:"customer_number"`
	CustomerName     string `json:"customer_name,omitempty"`
	BankCode         string `json:"bank_code"`
	BankName         string `json:"bank_name"`
	Remark           string `json:"remark"`
	OccupationSender string `json:"sender_job"`
	SenderCityName   string `json:"sender_city"`
}
type respData struct {
	BankName       string `json:"bank_name,omitempty"`
	BankCode       string `json:"bank_code,omitempty"`
	CustomerNumber string `json:"beneficiary_number,omitempty"`
	CustomerName   string `json:"beneficiary_name,omitempty"`
}

func (serv *DomesticServiceImpl) CreatePaycode(ctx *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[request]](ctx.Body())
	if errData != nil {
		return model.Response{
			Status: fiber.StatusBadRequest,
			RC:     rc.FAILED.String(),
			Errors: errData,
		}
	}
	masterStore, err := serv.DomesticRepository.PartnerMasterStoreId(req.Request.MasterStoreId)
	if err != nil {
		return model.Response{
			Status: err.Status,
			RC:     err.RC,
			Errors: append([]model.ErrorData{}, *err),
		}
	}
	senderData, errss := serv.DomesticRepository.DataAccountAll(req.Request.PhonePrefix, req.Request.Phone)
	if errss != nil {
		return model.Response{
			Status: errss.Status,
			RC:     errss.RC,
			Errors: append([]model.ErrorData{}, *errss),
		}
	}

	sopData, errss := serv.DomesticRepository.DataInitiator(req.Request.SopId)
	if errss != nil {
		return model.Response{
			Status: errss.Status,
			RC:     errss.RC,
			Errors: append([]model.ErrorData{}, *errss),
		}
	}
	imgSelf, _ := utils.GetImgString(*senderData.ImgSelf)
	imgSign, _ := utils.GetImgString(*senderData.ImgSign)
	ImgIdentity, _ := utils.GetImgString(*senderData.ImgIdentity)

	initKey := os.Getenv("CRED_INIT_KEY")
	initName := os.Getenv("CRED_INIT_NAME")

	sendDataas := sendCreatePaycode{
		Signature:   utils.SignMD5(req.Extref + utils.SignSHA256(initKey) + initName),
		PartnerId:   masterStore.PartnerId,
		Extref:      req.Extref,
		Amount:      req.Request.Amount,
		Phone:       req.Request.PhonePrefix + req.Request.Phone,
		Data1:       req.Request.Additional.Remark,
		Data2:       "",
		ImgSelf:     imgSelf,
		ImgSign:     imgSign,
		ImgIdentity: ImgIdentity,
		Sender: sendSender{
			Name:           senderData.AccountName,
			Phone:          req.Request.PhonePrefix + req.Request.Phone,
			Address:        *senderData.Address,
			OccupationId:   *senderData.OccupationId,
			SopId:          req.Request.SopId,
			Gender:         *senderData.Gender,
			City:           *senderData.CityId,
			PurposeId:      req.Request.Sender.PurposeId,
			SofundingId:    req.Request.Sender.SofId,
			DOB:            senderData.DOB.Format("2006-01-02"),
			POB:            *senderData.POB,
			IdentityTypeId: *senderData.IdentityTypeId,
			IdentityNumber: *senderData.IdentityNumber,
			IdentityExp:    "2500-01-02",
			ProvinceName:   "",
			AdminId:        senderData.AdminId,
		},
		Receiver: sendReceiver{
			Phone:          req.Request.Receiver.PhonePrefix + req.Request.Receiver.Phone,
			City:           req.Request.Receiver.CityId,
			Address:        "",
			CustomerName:   &req.Request.Additional.BenerName,
			IdentityTypeID: "",
			IdentityNumber: "",
		},
		Additional: sendAdditional{
			CustomerNumber:   req.Request.Additional.BenerNumber,
			CustomerName:     req.Request.Additional.BenerName,
			BankCode:         masterStore.BankCode,
			BankName:         masterStore.BankName,
			OccupationSender: "employee",
			Remark:           req.Request.Additional.Remark,
			SenderCityName:   *senderData.CityName,
		},
	}

	url := *sopData.InitiatorAdapAddr + "/create_paycode_sendmoney"
	statusCode, mapBody, errResponse := utils.PostSendToUrl(sendDataas, url)
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
			Status:  fiber.StatusOK,
			RC:      rc.FAILED.String(),
			Message: rc.FAILED.Message(),
			Extref:  req.Extref,
		}
	}

	return model.Response{
		Status:  fiber.StatusOK,
		RC:      rc.SUCCESS.String(),
		Message: rc.SUCCESS.Message(),
		Extref:  req.Extref,
		Data: resultResCrePay{
			VoucherExp:  mapBody["voucher_expire"].(string),
			VoucherCode: mapBody["voucher"].(string),
			Journey:     mapBody["journey"].(string),
		},
	}
}

type sendCreatePaycode struct {
	Signature   string         `json:"signature"`
	Extref      string         `json:"ext_ref"`
	Amount      string         `json:"transaction_amount"`
	Phone       string         `json:"customer_phone_number"`
	PartnerId   string         `json:"partner_id"`
	Data1       string         `json:"data_1"`
	Data2       string         `json:"data_2"`
	ImgSelf     string         `json:"selfie_picture"`
	ImgSign     string         `json:"signature_picture"`
	ImgIdentity string         `json:"identity_picture"`
	Sender      sendSender     `json:"data_sender"`
	Receiver    sendReceiver   `json:"data_receiver"`
	Additional  sendAdditional `json:"additional_data"`
}
type resultResCrePay struct {
	VoucherExp  string `json:"voucher_expired"`
	VoucherCode string `json:"voucher_code"`
	Journey     string `json:"journey"`
}
