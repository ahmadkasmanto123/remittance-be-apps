package model

// for data from database only
type CalculotorPayment struct {
	SopName     string `json:"sop_name"`
	InitiatorId int    `json:"initiator_id"`
}

type Destination struct {
	CountryId        string  `json:"country_id"`
	CountryName      string  `json:"country_name"`
	CountryCodeAlpha *string `json:"country_code_alpha"`
}

type Transaction struct {
	TransactionType string `json:"transaction_type"`
	MasterStoreId   string `json:"master_store_id"`
	BankName        string `json:"bank_name"`
}

type MasterStore struct {
	CountryId string `json:"country_id"`
	PartnerId string `json:"partner_id"`
	BankName  string `json:"bank_name"`
	BankCode  string `json:"bank_code"`
}
