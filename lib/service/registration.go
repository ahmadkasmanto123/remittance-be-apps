package service

import (
	"fmt"
	"love-remittance-be-apps/core/config"
	"love-remittance-be-apps/core/rc"
	"love-remittance-be-apps/core/utils"
	"love-remittance-be-apps/lib/interfc"
	"love-remittance-be-apps/lib/model"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type RegistrationServiceImpl struct {
	RegistrationRepository interfc.RegistrationRepository
}

func NewRegistrationService(repository interfc.RegistrationRepository) interfc.RegistrationService {
	return &RegistrationServiceImpl{
		RegistrationRepository: repository,
	}
}

func (serv *RegistrationServiceImpl) NewCustomer(ctx *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[model.CreateCustomer]](ctx.Body())

	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Extref:  req.Extref,
			Errors:  errData,
			Message: "Missing parameter",
		}
	}
	fmt.Println("check sini 1")
	// check email has been used before
	countEmail, errEmail := serv.RegistrationRepository.CountEmail(req.Request.Email)
	if errEmail != nil {
		return model.Response{
			Status: errEmail.Status,
			RC:     errEmail.RC,
			Extref: req.Extref,
			Errors: append([]model.ErrorData{}, *errEmail),
		}
	}
	if countEmail.Count != 0 {
		rcMessage1, errMsg := rc.GetResponseMsg("0006", rc.Register.Id(), req.Lang)
		if errMsg != nil {
			return model.Response{
				Status: fiber.StatusInternalServerError,
				RC:     errMsg.RC,
				Extref: req.Extref,
				Errors: append([]model.ErrorData{}, *errMsg),
			}
		}
		return model.Response{
			Status:  rcMessage1.Status,
			RC:      rcMessage1.RC,
			Extref:  req.Extref,
			Message: rcMessage1.Message,
		}
	}
	//check account
	countMsAccount, errMsCount := serv.RegistrationRepository.CountMsAccount(req.Request.PhonePrefix + req.Request.Phone)
	if errMsCount != nil {
		return model.Response{
			Status: errMsCount.Status,
			RC:     errMsCount.RC,
			Extref: req.Extref,
			Errors: append([]model.ErrorData{}, *errMsCount),
		}
	}
	fmt.Print("count ms account")
	fmt.Println(countMsAccount.Count)
	msAccount := false
	if countMsAccount.Count > 0 {
		countAccount, errAccount := serv.RegistrationRepository.CountAccount(req.Request.Phone, req.Request.PhonePrefix)
		fmt.Print("count ms admin")
		fmt.Println(countMsAccount.Count)
		if errAccount != nil {
			return model.Response{
				Status: errAccount.Status,
				RC:     errAccount.RC,
				Extref: req.Extref,
				Errors: append([]model.ErrorData{}, *errAccount),
			}
		}
		if countAccount.Count > 0 {
			fmt.Print("count ms accounttsdin")
			fmt.Println(countMsAccount.Count)
			rcMessage1, errMsg := rc.GetResponseMsg("0007", rc.Register.Id(), req.Lang)
			if errMsg != nil {
				return model.Response{
					Status: fiber.StatusInternalServerError,
					RC:     errMsg.RC,
					Extref: req.Extref,
					Errors: append([]model.ErrorData{}, *errMsg),
				}
			}
			return model.Response{
				Status:  rcMessage1.Status,
				RC:      rcMessage1.RC,
				Extref:  req.Extref,
				Message: rcMessage1.Message,
			}
		} else {
			fmt.Print("boolena ms account")
			fmt.Println(msAccount)
			msAccount = true
		}
	}

	fmt.Print("boolena2 ms account")
	fmt.Println(msAccount)
	storeCheck, errStore := serv.RegistrationRepository.Storecheck(req.Request.PhonePrefix)
	if errStore != nil {
		return model.Response{
			Status: fiber.StatusInternalServerError,
			RC:     errStore.RC,
			Extref: req.Extref,
			Errors: append([]model.ErrorData{}, *errStore),
		}
	}
	fmt.Print("storeCheck= ")
	fmt.Println(*storeCheck)

	var accountId *int
	if msAccount {
		//update ms account
		var errUpdate *model.ErrorData
		accountId, errUpdate = serv.RegistrationRepository.UpdateMsAccount(req.Request.FirstName+" "+req.Request.LastName, req.Request.PhonePrefix+req.Request.Phone)
		if errUpdate != nil {
			return model.Response{
				Status: fiber.StatusInternalServerError,
				RC:     errUpdate.RC,
				Extref: req.Extref,
				Errors: append([]model.ErrorData{}, *errUpdate),
			}
		}
		fmt.Print("accountId update= ")
		fmt.Println(*accountId)
	} else {
		var errInsert *model.ErrorData
		fmt.Println("masuk sinidfsdnf9")
		accountId, errInsert = serv.RegistrationRepository.InsertMsAccount(req.Request.FirstName+" "+req.Request.LastName, req.Request.PhonePrefix+req.Request.Phone)
		if errInsert != nil {
			return model.Response{
				Status: fiber.StatusInternalServerError,
				RC:     errInsert.RC,
				Extref: req.Extref,
				Errors: append([]model.ErrorData{}, *errInsert),
			}
		}
		if accountId == nil {
			return model.Response{
				Status:  fiber.StatusInternalServerError,
				RC:      "1001",
				Extref:  req.Extref,
				Message: "error insert",
			}
		}
		fmt.Print("accountId insert=  ")
		fmt.Println(*accountId)
	}

	adminId, errAdmin := serv.RegistrationRepository.InsertMsStoreAdmin(*req, storeCheck.StoreId, *accountId)
	if errAdmin != nil {
		return model.Response{
			Status: fiber.StatusInternalServerError,
			RC:     errAdmin.RC,
			Extref: req.Extref,
			Errors: append([]model.ErrorData{}, *errAdmin),
		}
	}
	if *adminId == 0 {
		return model.Response{
			Status:  fiber.StatusInternalServerError,
			RC:      "1001",
			Extref:  req.Extref,
			Message: "error insert",
		}
	}
	//generate otp
	otp := utils.GenerateOtp()
	//create data too session
	utils.SetSession(req.Request.PhonePrefix + req.Request.Phone + "_RegOtp")
	utils.AddSession(req.Request.PhonePrefix+req.Request.Phone+"_RegOtp", "regOtp", otp)

	fmt.Print("adminId= ")
	fmt.Println(*adminId)

	rcMessage1, errMsg := rc.GetResponseMsg("0000", rc.Register.Id(), req.Lang)
	if errMsg != nil {
		return model.Response{
			Status: fiber.StatusInternalServerError,
			RC:     errMsg.RC,
			Extref: req.Extref,
			Errors: append([]model.ErrorData{}, *errMsg),
		}
	}
	return model.Response{
		Status:  rcMessage1.Status,
		RC:      rcMessage1.RC,
		Extref:  req.Extref,
		Message: rcMessage1.Message,
		Data:    otp,
	}
}

type validateOtp struct {
	Phone       string `validate:"required" json:"phone"`
	PhonePrefix string `validate:"required" json:"phone_prefix"`
	Otp         int    `validate:"required" json:"otp"`
}

func (serv *RegistrationServiceImpl) RegOtpValidate(ctx *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[validateOtp]](ctx.Body())

	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "Missing parameter",
		}
	}
	//check session
	objeks := utils.GetSession(req.Request.PhonePrefix + req.Request.Phone + "_RegOtp")
	value, ok := objeks["regOtp"]
	var otp int
	if ok {
		otp, _ = strconv.Atoi(value.(string))
	} else {
		rcMessage1, errMsg := rc.GetResponseMsg("0003", rc.Login.Id(), req.Lang)
		if errMsg != nil {
			return model.Response{
				Status: fiber.StatusInternalServerError,
				RC:     errMsg.RC,
				Errors: append([]model.ErrorData{}, *errMsg),
			}
		}
		return model.Response{
			Status:  rcMessage1.Status,
			RC:      rcMessage1.RC,
			Message: rcMessage1.Message,
		}
	}
	//check data account
	count, errcount := serv.RegistrationRepository.CountAccount(req.Request.Phone, req.Request.PhonePrefix)
	if errcount != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  append([]model.ErrorData{}, *errcount),
			Message: "database",
		}
	}

	if count.Count < 1 {
		rcMessage1, errMsg := rc.GetResponseMsg("0003", rc.Login.Id(), req.Lang)
		if errMsg != nil {
			return model.Response{
				Status: fiber.StatusInternalServerError,
				RC:     errMsg.RC,
				Errors: append([]model.ErrorData{}, *errMsg),
			}
		}
		return model.Response{
			Status:  rcMessage1.Status,
			RC:      rcMessage1.RC,
			Message: rcMessage1.Message,
		}
	}
	if otp != req.Request.Otp {
		rcMessage1, errMsg := rc.GetResponseMsg("0005", rc.Login.Id(), req.Lang)
		if errMsg != nil {
			return model.Response{
				Status: fiber.StatusInternalServerError,
				RC:     errMsg.RC,
				Errors: append([]model.ErrorData{}, *errMsg),
			}
		}
		return model.Response{
			Status:  rcMessage1.Status,
			RC:      rcMessage1.RC,
			Message: rcMessage1.Message,
		}
	}
	//get data account
	dataAccount, errAccount := serv.RegistrationRepository.DataAccountReg(req.Request.Phone, req.Request.PhonePrefix)
	if errAccount != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  append([]model.ErrorData{}, *errAccount),
			Message: "database",
		}
	}

	//update ms_store_admin
	err := serv.RegistrationRepository.UpdateMsStoreAdmin(req.DeviceId, dataAccount.AdminId)
	if err != nil {
		return model.Response{
			Status: fiber.StatusInternalServerError,
			RC:     err.RC,
			Errors: append([]model.ErrorData{}, *err),
		}
	}

	//delete session
	utils.DelDataSession(req.Request.PhonePrefix + req.Request.Phone + "_RegOtp")

	token, errss := config.CreateToken(dataAccount.AccountName, req.Request.Phone, req.Request.PhonePrefix)
	if errss != nil {
		return model.Response{
			Status:  fiber.StatusUnauthorized,
			RC:      rc.UNAUTHORIZED.String(),
			Message: "Token not created",
		}
	}
	return model.Response{
		Status:  fiber.StatusOK,
		RC:      rc.SUCCESS.String(),
		Message: rc.SUCCESS.Message(),
		Extref:  req.Extref,
		Data: model.LoginResponse{
			Token: token,
		},
	}
}
