package utils

import (
	"log"
	"love-remittance-be-apps/core/rc"
	"love-remittance-be-apps/lib/model"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jlaffaye/ftp"
)

func PostSendToUrl(T any, url string) (int, map[string]interface{}, *model.ErrorData) {
	agent := fiber.Post(url)
	agent.JSON(T)
	agent.InsecureSkipVerify()
	log.Printf("ToUrl = %s | Method = %s \n| RequestBody = %s", url, agent.Request().Header.Method(), string(agent.Request().Body()))
	statusCode, body, errResponse := agent.Bytes()
	log.Printf("FromUrl = %s | StatusCode = %d | Response = %s", url, statusCode, string(body))
	if len(errResponse) > 0 {
		return statusCode, ToMap(body), &model.ErrorData{
			Status:      statusCode,
			RC:          rc.FAILED.String(),
			Description: "No Response from client",
		}
	}

	return statusCode, ToMap(body), nil
}

func PostForSession(T any, url string) (int, map[string]interface{}, []error) {
	agent := fiber.Post(url)
	agent.JSON(T)
	agent.InsecureSkipVerify()
	log.Printf("ToUrl = %s | Method = %s \n| RequestBody = %s", url, agent.Request().Header.Method(), string(agent.Request().Body()))
	statusCode, body, errResponse := agent.Bytes()
	log.Printf("FromUrl = %s | StatusCode = %d | Response = %s", url, statusCode, string(body))
	if len(errResponse) > 0 {
		return statusCode, ToMap(body), errResponse
	}

	return statusCode, ToMap(body), nil
}

func FtpConnection() (*ftp.ServerConn, *model.ErrorData) {
	conns, err := ftp.Dial(os.Getenv("FTP_ADDRESS"))
	if err != nil {
		return nil, &model.ErrorData{
			Description: "Error connect to server ftp " + err.Error(),
		}
	}

	err = conns.Login(os.Getenv("FTP_USER"), os.Getenv("FTP_PASSWORD"))
	if err != nil {
		return nil, &model.ErrorData{
			Description: "Error login to server ftp " + err.Error(),
		}
	}
	return conns, nil
}
