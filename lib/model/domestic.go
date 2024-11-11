package model

type DataSOPayment struct {
	SopType           *string `json:"type,omitempty"`
	SopName           *string `json:"sop_name"`
	SopId             *int    `json:"sop_id,omitempty"`
	InitiatorId       *int    `json:"initiator_id"`
	InitiatorAdapAddr *string `json:"initiator_adapter_address,omitempty"`
}

type DynamicData[T any] struct {
	DataType  string `json:"type"`
	DataItems T      `json:"items"`
}

type DataDestination struct {
	Types         string  `json:"type,omitempty"`
	BankName      string  `json:"bank_name"`
	BankCode      *string `json:"bank_code,omitempty"`
	MasterStoreId int     `json:"master_store_id"`
}
