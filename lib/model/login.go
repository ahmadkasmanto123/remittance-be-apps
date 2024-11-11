package model

type CountSomething struct {
	Count int `json:"count"`
}

type DataAccount struct {
	AccountId       int    `json:"account_id"`
	AccountStatutId int    `json:"account_status_id"`
	AccountName     string `json:"account_name"`
	AdminEmail      string `json:"admin_email"`
	AdminPin        string `json:"admin_pin"`
	AdminConf       string `json:"admin_conf"`
	DeviceId        string `json:"device_id"`
	AdminId         int    `json:"admin_id"`
}
