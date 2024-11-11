package utils

import (
	"os"
	"strconv"
)

type SesionReq struct {
	SessionId  string  `json:"sessId"`
	SessionTtl *string `json:"sessionTtl,omitempty"` //time in second
}

var ttlDefault = "3600"

func SetSession(name string) map[string]interface{} {
	request := SesionReq{
		SessionId:  name,
		SessionTtl: &ttlDefault,
	}
	newUrl := os.Getenv("SESSION_URL") + "/setSession"
	_, body, _ := PostForSession(request, newUrl)
	return body
}

func RenewSession(name string) map[string]interface{} {
	request := SesionReq{
		SessionId:  name,
		SessionTtl: &ttlDefault,
	}
	newUrl := os.Getenv("SESSION_URL") + "/renewSession"
	_, body, _ := PostForSession(request, newUrl)
	return body
}

type DataSession struct {
	SessionId string `json:"sessId"`
	Key       string `json:"key"`
	Value     any    `json:"value"`
}

func AddSession(name string, nameKey string, T any) map[string]interface{} {
	request := DataSession{
		SessionId: name,
		Key:       nameKey,
		Value:     T,
	}
	newUrl := os.Getenv("SESSION_URL") + "/addDataSession"
	_, body, _ := PostForSession(request, newUrl)
	return body
}

func GetSession(name string) map[string]interface{} {
	request := SesionReq{
		SessionId: name,
	}
	newUrl := os.Getenv("SESSION_URL") + "/getSession"
	_, body, _ := PostForSession(request, newUrl)
	return body
}

func CheckExistSession(name string) bool {
	request := SesionReq{
		SessionId: name,
	}
	newUrl := os.Getenv("SESSION_URL") + "/checkExistCookie"
	_, body, _ := PostForSession(request, newUrl)
	result, _ := strconv.ParseBool(body["result"].(string))
	return result
}

func CheckExpireSession(name string) bool {
	request := SesionReq{
		SessionId: name,
	}
	newUrl := os.Getenv("SESSION_URL") + "/checkExpireCookie"
	_, body, _ := PostForSession(request, newUrl)
	result, _ := strconv.ParseBool(body["result"].(string))
	return result
}

func DelDataSession(name string) bool {
	request := SesionReq{
		SessionId: name,
	}
	newUrl := os.Getenv("SESSION_URL") + "/delDataSession"
	_, body, _ := PostForSession(request, newUrl)
	result := false
	if body["pesan"].(string) == "OK" {
		result = true
	}
	return result
}
