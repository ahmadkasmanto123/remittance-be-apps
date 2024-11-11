package model

type Param struct {
	Limit  *int    `json:"limit,omitempty"`
	Offset *int    `json:"offset,omitempty"`
	Search *string `json:"search,omitempty"`
	Sort   *string `json:"sort,omitempty"`
}

type DefaultRequest[T any] struct {
	// Param   Param  `json:"param,omitempty"`
	Request  T      `json:"data"`
	Extref   string `validate:"required" json:"extref"`
	Lang     string `validate:"required" json:"lang"`
	DeviceId string `validate:"required" json:"device_id"`
}
