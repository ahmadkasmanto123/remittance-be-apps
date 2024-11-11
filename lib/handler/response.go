package handler

type Response struct {
	RC            string      `json:"rc"`
	Message       string      `json:"message"`
	Detail        string      `json:"detail"`
	TransactionId string      `json:"transaction_id"`
	Extref        string      `json:"ext_ref"`
	Data          interface{} `json:"data"`
}
