package model

type StoreCheck struct {
	StoreId   int `json:"store_id"`
	CountryId int `json:"country_id"`
}

type CreateCustomer struct {
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Phone       string `json:"phone" validate:"required"`
	PhonePrefix string `json:"phone_prefix" validate:"required"`
	Pin         string `json:"pin" validate:"required"`
}
