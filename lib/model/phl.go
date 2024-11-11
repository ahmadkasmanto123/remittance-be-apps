package model

type ReqPhl struct {
	Signature   string           `json:"signature"`
	Extref      string           `json:"ext_ref"`
	Amount      string           `json:"transaction_amount"`
	Phone       string           `json:"customer_phone_number"`
	PartnerId   int              `json:"partner_id"`
	Data1       string           `json:"data_1"`
	Data2       string           `json:"data_2"`
	ImgSelf     string           `json:"selfie_picture"`
	ImgSign     string           `json:"signature_picture"`
	ImgIdentity string           `json:"identity_picture"`
	Sender      SenderReqPhl     `json:"data_sender"`
	Receiver    ReceiverReqPhl   `json:"data_receiver"`
	Additional  AdditionalReqPhl `json:"additional_data"`
}
type SenderReqPhl struct {
	Name           string `json:"name"`
	Phone          string `json:"phone_number"`
	Address        string `json:"address"`
	OccupationId   int    `json:"occupation"`
	SopId          int    `json:"source_of_payment"`
	Gender         string `json:"gender"`
	PurposeId      int    `json:"purpose"`
	SofundingId    int    `json:"source_of_fund"`
	DOB            string `json:"dob"`
	POB            string `json:"pob"`
	IdentityTypeId int    `json:"identity_type"`
	IdentityNumber string `json:"identity_number"`
	IdentityExp    string `json:"identity_exp"`
	AdminId        string `json:"admin_id"`
	City           int    `json:"city"`
	Data1          string `json:"data_1"`
}
type ReceiverReqPhl struct {
	FirstName      string `json:"firstname"`
	LastName       string `json:"lastname"`
	Address        string `json:"address"`
	City           int    `json:"city"`
	Name           string `json:"name"`
	Phone          string `json:"phone_number"`
	IdentityNumber string `json:"identity_number"`
	IdentityTypeID int    `json:"identity_type"`
}
type AdditionalReqPhl struct {
	FirstName            string `json:"firstname"`
	LastName             string `json:"lastname"`
	States               string `json:"state_receiver"`
	PostalCodeReceiver   string `json:"postal_code_receiver"`
	IdentityTypeReceiver string `json:"identity_type_receiver"`
	PhoneReceiver        string `json:"phone_number_receiver"`
	CityReceiver         string `json:"city_name_receiver"`
	AddressReceiver      string `json:"address_receiver"`
	NationalityReceiver  string `json:"receiver_nationality"`
	CustomerNumber       string `json:"customer_number"`
	CustomerName         string `json:"customer_name,omitempty"`
	BankBranchCode       string `json:"bank_branch_code"`
	BankCode             string `json:"bank_code"`
	BankName             string `json:"bank_name"`
	Remark               string `json:"remark"`
	PostalCodeSender     string `json:"postal_code_sender"`
	OccupationSender     string `json:"occupation_partner_key"`
	SofundingSender      string `json:"sof_partner_key"`
	PurposeSender        string `json:"purpose_partner_key"`
	RelationSender       string `json:"relationship_partner_key"`
	CitySender           string `json:"sender_city"`
	AmountReceiver       string `json:"amount_receiver"`
}
