package service

import (
	"love-remittance-be-apps/core/rc"
	"love-remittance-be-apps/core/utils"
	"love-remittance-be-apps/lib/interfc"
	"love-remittance-be-apps/lib/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UpdateServiceImpl struct {
	UpdateRepository interfc.UpdateRepository
}

func NewUpdateService(repository interfc.UpdateRepository) interfc.UpdateService {
	return &UpdateServiceImpl{
		UpdateRepository: repository,
	}
}

func (serv *UpdateServiceImpl) Profile(c *fiber.Ctx) model.Response {
	req, errData := utils.JsonToObject[model.DefaultRequest[model.UpdateProfil]]([]byte(c.FormValue("request")))
	if errData != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  errData,
			Message: "error sini",
			Extref:  req.Extref,
		}
	}
	//check login account
	countAccount, errcount := serv.UpdateRepository.CountAccount(req.Request.Phone, req.Request.PhonePrefix)
	if errcount != nil {
		return model.Response{
			Status:  fiber.StatusBadRequest,
			RC:      rc.FAILED.String(),
			Errors:  append([]model.ErrorData{}, *errcount),
			Message: "database",
			Extref:  req.Extref,
		}
	}
	if countAccount.Count == 0 {
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

	//upload image to ftp
	uuId := uuid.New().String()
	nameFile := []string{"img_signature", "img_self", "img_self_identity"}
	var resultName []string
	for i := 0; i < len(nameFile); i++ {
		rsltName, err := utils.UploadImg(c, nameFile[i], uuId)
		if err != nil {
			return model.Response{
				RC:      rc.INTERNALSERVERERROR.String(),
				Status:  fiber.StatusInternalServerError,
				Message: rc.INTERNALSERVERERROR.Message(),
				Errors:  append([]model.ErrorData{}, *err),
				Extref:  req.Extref,
			}
		}
		resultName = append(resultName, *rsltName)
	}

	up, errUpdate := serv.UpdateRepository.UpdateMsAccount((*req), resultName)

	if errUpdate != nil {
		return model.Response{
			RC:      rc.INTERNALSERVERERROR.String(),
			Status:  fiber.StatusInternalServerError,
			Message: rc.INTERNALSERVERERROR.Message(),
			Errors:  append([]model.ErrorData{}, *errUpdate),
			Extref:  req.Extref,
		}
	}
	return model.Response{
		Status:  fiber.StatusOK,
		RC:      rc.SUCCESS.String(),
		Message: rc.SUCCESS.Message(),
		Extref:  req.Extref,
		Data:    up,
	}
}
