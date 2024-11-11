package model

import "time"

type ShowImgData struct {
	Phone       string `json:"phone" validate:"required"`
	PhonePrefix string `json:"phone_prefix" validate:"required"`
	ImgName     string `json:"img_name" validate:"required"`
}

type Province struct {
	ProvinceId   int    `json:"province_id"`
	ProvinceName string `json:"province_name,omitempty"`
}

type DataIdName struct {
	Id   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

type DataAccountDetail struct {
	AccountId         int        `json:"id"`
	AccountStatutId   int        `json:"status_id"`
	AccountStatusName string     `json:"status_name"`
	IdentityTypeId    *int       `json:"identity_type_id,omitempty"`
	IdentityTypeName  *string    `json:"identity_type_name,omitempty"`
	IdentityNumber    *string    `json:"identity_number,omitempty"`
	AccountName       string     `json:"account_name"`
	AdminId           int        `json:"admin_id"`
	AdminEmail        string     `json:"admin_email"`
	DeviceId          string     `json:"device_id"`
	CityId            *int       `json:"city_id,omitempty"`
	CityName          *string    `json:"city_name,omitempty"`
	OccupationId      *int       `json:"occupation_id,omitempty"`
	OccupationName    *string    `json:"occupation_name,omitempty"`
	Address           *string    `json:"address,omitempty"`
	POB               *string    `json:"pob,omitempty"`
	DOB               *time.Time `json:"dob,omitempty"`
	Gender            *string    `json:"gender,omitempty"`
	PostalCode        *string    `json:"postal_code,omitempty"`
	ImgSelf           *string    `json:"img_self,omitempty"`
	ImgIdentity       *string    `json:"img_identity,omitempty"`
	ImgSign           *string    `json:"img_sign,omitempty"`
}

type History struct {
	TransactionId     int    `json:"transaction_id"`
	ExtRef            string `json:"ext_ref"`
	TransactionAmount string `json:"transaction_amount"`
	TimeTrx           string `json:"time_trx"`
	Region            string `json:"region"`
	BankName          string `json:"bank_name"`
	ReceiverCurrency  string `json:"receiver_currency"`
	SenderCurrency    string `json:"sender_currency"`
	Status            int    `json:"status"`
	TypeTrx           string `json:"type_trx"`
	TypeId            int    `json:"type_id,omitempty"`
	UrlAdapter        string `json:"url_adapter,omitempty"`
	ReceiverName      string `json:"receiver_name,omitempty"`
	ReceiverPhone     string `json:"receiver_phone,omitempty"`
	ReceiverCountry   string `json:"receiver_country,omitempty"`
}

type DataPhone struct {
	Phone       string `json:"phone" validate:"required"`
	PhonePrefix string `json:"phone_prefix" validate:"required"`
}

type VoucherActive struct {
	BankName           string `json:"bank_name"`
	BankBranchCode     string `json:"bank_branch_code,omitempty"`
	BenefiaciaryName   string `json:"beneficiary_name"`
	BenefiaciaryNumber string `json:"beneficiary_number"`
	Remark             string `json:"remark"`
	Amount             string `json:"amount"`
	AmountFee          string `json:"amount_fee"`
	AmountTotal        string `json:"amount_total"`
	AmountReceive      string `json:"amount_receive,omitempty"`
	PaycodeNumber      string `json:"paycode_number"`
	Journey            string `json:"journey"`
	TitlePayment       string `json:"title_payment"`
}

type HistoryDetail struct {
	BankName           string `json:"bank_name"`
	BankBranchCode     string `json:"bank_branch_code,omitempty"`
	BenefiaciaryName   string `json:"beneficiary_name"`
	BenefiaciaryNumber string `json:"beneficiary_number"`
	Remark             string `json:"remark"`
	Amount             string `json:"amount"`
	AmountFee          string `json:"amount_fee"`
	AmountTotal        string `json:"amount_total"`
	AmountReceive      string `json:"amount_receive,omitempty"`
	Rate               string `json:"rate"`
	SenderName         string `json:"sender_name"`
	SenderAddress      string `json:"sender_address"`
	TransactionId      int    `json:"transaction_id"`
	DateTrx            string `json:"date_trx"`
	TypeTrx            string `json:"type_trx"`
	FootNote           string `json:"foot_note"`
	ReceiverName       string `json:"receiver_name,omitempty"`
	ReceiverPhone      string `json:"receiver_phone,omitempty"`
	ReceiverCountry    string `json:"receiver_country,omitempty"`
}
