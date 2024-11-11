package rc

type RC int

const (
	SUCCESS RC = iota
	FAILED
	PENDING
	NOTFOUND
	UNAUTHORIZED
	EXTREFDOUBLE
	INTERNALSERVERERROR
)

func (rc RC) Message() string {
	return [...]string{
		"Transaction Success",
		"Transaction Failed",
		"Transaction Pending",
		"Not Found",
		"Unauthorized",
		"Extref has been used",
		"Internal server Error",
	}[rc]
}

func (rc RC) String() string {
	return [...]string{
		"0000",
		"1001",
		"1002",
		"1003",
		"1004",
		"1005",
		"1006",
	}[rc]
}

type MSTYPE int

const (
	Common MSTYPE = iota
	Register
	Login
	Billing
	Invoice
	OtpLoginCreate
	OtpLoginCheck
	Logout
	Customer
	Company
	Company_KYB
	Company_Verify
	Account
	Account_Verify
	Account_Profile
	Account_Change_Email
)

func (msType MSTYPE) Id() int {
	return [...]int{
		1,
		2,
		3,
		4,
		5,
		6,
		7,
		8,
		9,
		10,
		11,
		12,
		13,
		14,
		15,
		16,
	}[msType]
}
