package models

type Employee struct {
	ID           int64  `json:"id"`
	OrgId        int64  `json:"org_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	MobileNumber string `json:"mobile_number"`
}
