package model

type CountryData struct {
	CountryId   int    `json:"country_id"`
	CountryName string `json:"country_name"`
	Currency    string `json:"currency"`
	Code2       string `json:"code_2"`
	Code3       string `json:"code_3"`
}

type BankDest struct {
	MasterStoreId int    `json:"master_store_id"`
	BankName      string `json:"bank_name"`
	TypeDest      string `json:"type"`
}

type Credential struct {
	BankName      string `json:"bank_name"`
	BankCode      string `json:"bank_code"`
	PartnerId     int    `json:"partner_id"`
	CountryId     int    `json:"country_id"`
	InitiatorId   int    `json:"initiator_id"`
	InitiatorUser string `json:"initiator_user"`
	InitiatorKey  string `json:"initiator_key"`
	AdminId       string `json:"admin_id"`
	CountryCode   string `json:"country_code"`
	UrlAdapter    string `json:"url_adapter"`
	CityReceiver  int    `json:"city"`
}

type Mandatory struct {
	SopId         int            `validate:"required" json:"sop_id"`
	MasterStoreId int            `validate:"required" json:"master_store_id"`
	Amount        string         `validate:"required" json:"transaction_amount"`
	AmountReceive string         `validate:"required" json:"transaction_amount_receiver"`
	PhonePrefix   string         `validate:"required" json:"phone_prefix"`
	Phone         string         `validate:"required" json:"phone"`
	Sender        SenderMand     `validate:"required" json:"sender"`
	Receiver      ReceiverMand   `validate:"required" json:"receiver"`
	Additional    AdditionalMand `validate:"required" json:"additional"`
}
type SenderMand struct {
	SofundingId  int `validate:"required" json:"sofunding_id"`
	PurposeId    int `validate:"required" json:"purpose_id"`
	OccupationId int `validate:"required" json:"occupation_id"`
	RelationId   int `validate:"required" json:"relation_id"`
}
type ReceiverMand struct {
	Phone          string `validate:"required" json:"phone"`
	IdentityType   int    `validate:"required" json:"identity_type"`
	IdentityNumber string `validate:"required" json:"identity_number"`
	FirstName      string `validate:"required" json:"first_name"`
	LastName       string `validate:"required" json:"last_name"`
	Address        string `validate:"required" json:"address"`
	CityName       string `validate:"required" json:"city_name"`
	States         string `json:"states,omitempty"`
	PostalCode     string `json:"postal_code,omitempty"`
	CountryId      int    `validate:"required" json:"country_id"`
}
type AdditionalMand struct {
	BenefiaciaryNumber string `validate:"required" json:"beneficiary_number"`
	BenefiaciaryName   string `validate:"required" json:"beneficiary_name"`
	Remark             string `json:"remark,omitempty"`
	BankBranchCode     string `json:"bank_branch_code,omitempty"`
}

type DataIdNkey struct {
	Id  int    `json:"id"`
	Key string `json:"key"`
}

type AllDataIdnKey struct {
	Identitas   DataIdNkey `json:"identitas"`
	Relations   DataIdNkey `json:"relations"`
	Purposes    DataIdNkey `json:"purposes"`
	Fundings    DataIdNkey `json:"fundings"`
	Occupations DataIdNkey `json:"occupations"`
}
