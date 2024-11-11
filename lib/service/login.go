package service

import (
	"love-remittance-be-apps/core/config"
	"love-remittance-be-apps/core/rc"
	"love-remittance-be-apps/core/utils"
	"love-remittance-be-apps/lib/interfc"
	"love-remittance-be-apps/lib/model"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type LoginServiceImpl struct {
	LoginRepository interfc.LoginRepository
}

func NewLoginService(repository interfc.LoginRepository) interfc.LoginService {
	return &LoginServiceImpl{
		LoginRepository: repository,
	}
}

type LoginReq struct {
	Phone       string `json:"phone" validate:"required"`
	PhonePrefix string `json:"phone_prefix" validate:"required"`
	Pin         string `json:"pin" validate:"required"`
}

func (serv *LoginServiceImpl) Login(ctx *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[LoginReq]](ctx.Body())

	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "Missing parameter",
			Extref:  req.Extref,
		}
	}
	count, error := serv.LoginRepository.CountAccount(req.Request.Phone, req.Request.PhonePrefix)
	if error != nil {
		return model.Response{
			Status:  400,
			RC:      "0001",
			Message: "Internal Server Error",
			Errors:  append([]model.ErrorData{}, *error),
			Extref:  req.Extref,
		}
	}

	if count.Count < 1 {
		rcMessage1, errMsg := rc.GetResponseMsg("0003", rc.Login.Id(), req.Lang)
		if errMsg != nil {
			return model.Response{
				Status: fiber.StatusInternalServerError,
				RC:     errMsg.RC,
				Errors: append([]model.ErrorData{}, *errMsg),
				Extref: req.Extref,
			}
		}
		return model.Response{
			Status:  rcMessage1.Status,
			RC:      rcMessage1.RC,
			Message: rcMessage1.Message,
			Extref:  req.Extref,
		}
	}

	// data, errDatas :=
	data, errDatas := serv.LoginRepository.DataAccount(req.Request.Phone, req.Request.PhonePrefix)
	if errDatas != nil {
		return model.Response{
			Status: fiber.StatusUnauthorized,
			RC:     errDatas.RC,
			Errors: append([]model.ErrorData{}, *errDatas),
			Extref: req.Extref,
		}
	}

	//check account status
	if data.AccountStatutId == 3 || data.AccountStatutId == 4 || data.AccountStatutId == 7 {
		rcMessage, errMsg := rc.GetResponseMsg("0003", rc.Login.Id(), req.Lang)
		if errMsg != nil {
			return model.Response{
				Status: fiber.StatusInternalServerError,
				RC:     errMsg.RC,
				Errors: append([]model.ErrorData{}, *errMsg),
				Extref: req.Extref,
			}
		}
		return model.Response{
			Status:  rcMessage.Status,
			RC:      rcMessage.RC,
			Message: "AKUN BLOK",
			Extref:  req.Extref,
		}
	} else if data.AccountStatutId == 8 {
		rcMessage, errMsg := rc.GetResponseMsg("0003", rc.Login.Id(), req.Lang)
		if errMsg != nil {
			return model.Response{
				Status: fiber.StatusInternalServerError,
				RC:     errMsg.RC,
				Errors: append([]model.ErrorData{}, *errMsg),
				Extref: req.Extref,
			}
		}
		return model.Response{
			Status:  rcMessage.Status,
			RC:      rcMessage.RC,
			Message: "AKUN CLOSED",
			Extref:  req.Extref,
		}
	}

	if data.AdminPin != req.Request.Pin {
		objeks := utils.GetSession(req.Request.PhonePrefix + req.Request.Phone + "_WrongPin")
		var counter int = 0
		rcMessage, errMsg := rc.GetResponseMsg("0003", rc.Login.Id(), req.Lang)
		message := rcMessage.Message
		if errMsg != nil {
			return model.Response{
				Status: fiber.StatusInternalServerError,
				RC:     errMsg.RC,
				Errors: append([]model.ErrorData{}, *errMsg),
				Extref: req.Extref,
			}
		}
		value, ok := objeks["wrongPin"]
		if ok {
			counter, _ = strconv.Atoi(value.(string))
			if counter == 1 {
				utils.AddSession(req.Request.PhonePrefix+req.Request.Phone+"_WrongPin", "wrongPin", 2)
				message = message + " " + strconv.Itoa(counter+1) + " kali"
			} else {
				err := serv.LoginRepository.UpdateAccountStatus(req.Request.PhonePrefix+req.Request.Phone, 4)
				if err != nil {
					return model.Response{
						Status: fiber.StatusInternalServerError,
						RC:     err.RC,
						Errors: append([]model.ErrorData{}, *err),
						Extref: req.Extref,
					}
				}
				message = message + " 3 kali akun sudah diblok"
			}
		} else {
			utils.SetSession(req.Request.PhonePrefix + req.Request.Phone + "_WrongPin")
			utils.AddSession(req.Request.PhonePrefix+req.Request.Phone+"_WrongPin", "wrongPin", 1)
			message = message + " 1 kali"
		}

		return model.Response{
			Status:  rcMessage.Status,
			RC:      rcMessage.RC,
			Message: message,
			Extref:  req.Extref,
		}
	}
	if data.DeviceId == "" {
		rcMessage3, errMsg := rc.GetResponseMsg("0004", rc.Login.Id(), req.Lang)
		if errMsg != nil {
			return model.Response{
				Status: fiber.StatusInternalServerError,
				RC:     errMsg.RC,
				Errors: append([]model.ErrorData{}, *errMsg),
				Extref: req.Extref,
			}
		}
		return model.Response{
			Status:  rcMessage3.Status,
			RC:      rcMessage3.RC,
			Message: rcMessage3.Message,
			Extref:  req.Extref,
		}
	}

	if data.DeviceId != req.DeviceId {
		rcMessage3, errMsg := rc.GetResponseMsg("0004", rc.Login.Id(), req.Lang)
		if errMsg != nil {
			return model.Response{
				Status: fiber.StatusInternalServerError,
				RC:     errMsg.RC,
				Errors: append([]model.ErrorData{}, *errMsg),
				Extref: req.Extref,
			}
		}
		return model.Response{
			Status:  rcMessage3.Status,
			RC:      rcMessage3.RC,
			Message: rcMessage3.Message,
			Extref:  req.Extref,
		}
	}

	// resultDevice := utils.StringToArray(data.AdminConf)
	// if len(resultDevice) == 0 {
	// 	rcMessage3, errMsg := rc.GetResponseMsg("0004", rc.Login.Id(), req.Lang)
	// 	if errMsg != nil {
	// 		return model.Response{
	// 			Status: fiber.StatusInternalServerError,
	// 			RC:     errMsg.RC,
	// 			Errors: append([]model.ErrorData{}, *errMsg),
	// 		}
	// 	}
	// 	return model.Response{
	// 		Status:  rcMessage3.Status,
	// 		RC:      rcMessage3.RC,
	// 		Message: rcMessage3.Message,
	// 	}
	// }
	// // Loop over string slice at key.
	// validDevice := false
	// for i := range resultDevice {
	// 	if req.DeviceId == resultDevice[i] {
	// 		validDevice = true
	// 		break
	// 	}
	// }
	// if !validDevice {
	// 	rcMessage3, errMsg := rc.GetResponseMsg("0004", rc.Login.Id(), req.Lang)
	// 	if errMsg != nil {
	// 		return model.Response{
	// 			Status: fiber.StatusInternalServerError,
	// 			RC:     errMsg.RC,
	// 			Errors: append([]model.ErrorData{}, *errMsg),
	// 		}
	// 	}
	// 	return model.Response{
	// 		Status:  rcMessage3.Status,
	// 		RC:      rcMessage3.RC,
	// 		Message: rcMessage3.Message,
	// 	}
	// }

	token, errss := config.CreateToken(data.AccountName, req.Request.Phone, req.Request.PhonePrefix)
	if errss != nil {
		return model.Response{
			Status:  fiber.StatusUnauthorized,
			RC:      rc.UNAUTHORIZED.String(),
			Message: "Token not created",
		}
	}

	utils.DelDataSession(req.Request.PhonePrefix + req.Request.Phone + "_WrongPin")
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

type createOtpReq struct {
	Phone       string `validate:"required" json:"phone"`
	PhonePrefix string `validate:"required" json:"phone_prefix"`
}

func (serv *LoginServiceImpl) LogInOtpCreate(ctx *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[createOtpReq]](ctx.Body())

	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "Missing parameter",
		}
	}
	//check data account
	count, errcount := serv.LoginRepository.CountAccount(req.Request.Phone, req.Request.PhonePrefix)
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
	//generate otp
	otp := utils.GenerateOtp()
	//create data too session
	utils.SetSession(req.Request.PhonePrefix + req.Request.Phone + "_LogOtp")
	utils.AddSession(req.Request.PhonePrefix+req.Request.Phone+"_LogOtp", "logOtp", otp)

	//send message with sms

	return model.Response{
		Status:  fiber.StatusOK,
		RC:      rc.SUCCESS.String(),
		Message: rc.SUCCESS.Message(),
		Extref:  req.Extref,
		Data:    otp,
	}
}

type validateOtpReq struct {
	Phone       string `validate:"required" json:"phone"`
	PhonePrefix string `validate:"required" json:"phone_prefix"`
	Otp         int    `validate:"required" json:"otp"`
}

func (serv *LoginServiceImpl) LogInOtpValidate(ctx *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[validateOtpReq]](ctx.Body())

	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "Missing parameter",
		}
	}
	//check session
	objeks := utils.GetSession(req.Request.PhonePrefix + req.Request.Phone + "_LogOtp")
	value, ok := objeks["logOtp"]
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
	count, errcount := serv.LoginRepository.CountAccount(req.Request.Phone, req.Request.PhonePrefix)
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
	dataAccount, errAccount := serv.LoginRepository.DataAccount(req.Request.Phone, req.Request.PhonePrefix)
	if errAccount != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  append([]model.ErrorData{}, *errAccount),
			Message: "database",
		}
	}

	//update ms_store_admin
	err := serv.LoginRepository.UpdateMsStoreAdmin(req.DeviceId, dataAccount.AdminId)
	if err != nil {
		return model.Response{
			Status: fiber.StatusInternalServerError,
			RC:     err.RC,
			Errors: append([]model.ErrorData{}, *err),
		}
	}

	//delete session
	utils.DelDataSession(req.Request.PhonePrefix + req.Request.Phone + "_LogOtp")

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
