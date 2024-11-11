package model

type UpdateProfil struct {
	Phone           string `validate:"required" json:"phone"`
	PhonePrefix     string `validate:"required" json:"phone_prefix"`
	FirstName       string `validate:"required" json:"first_name"`
	LastName        string `validate:"required" json:"last_name"`
	Email           string `validate:"required" json:"email"`
	IdentityTypeId  int    `validate:"required" json:"identity_type_id"`
	IdentityNumber  string `validate:"required" json:"identity_number"`
	IdentityExpired string `validate:"required" json:"identity_expired"`
	Address         string `validate:"required" json:"address"`
	City            int    `validate:"required" json:"city_id"`
	Occupation      int    `validate:"required" json:"occupation_id"`
	POB             string `validate:"required" json:"pob"`
	DOB             string `validate:"required" json:"dob"`
	PostalCode      string `validate:"required" json:"postal_code"`
	Gender          string `validate:"required" json:"gender"`
}
