package models

type Customer struct {
	FullName     string `json:"fullName"`
	MobileNumber string `json:"mobileNumber"`
	Email        string `json:"email"`
	CompanyName  string `json:"companyName"`
	Address_1    string `json:"address_1"`
	Address_2    string `json:"address_2"`
	City         string `json:"city"`
	District     string `json:"district"`
	State        string `json:"state"`
	Pincode      string `json:"pincode"`
	Type         string `json:"type"`
	Reference    string `json:"reference"`
	GstNumber    string `json:"gstNumber"`
}
