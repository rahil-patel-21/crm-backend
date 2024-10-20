// Packages
package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type Ticket struct {
	ID                int64     `json:"id"`
	Created_By        int64     `json:"created_by"`
	Customer_Name     string    `json:"customer_name"`
	Contact_Number    string    `json:"contact_number"`
	Email             string    `json:"email"`
	Address           string    `json:"address"`
	City              string    `json:"city"`
	District          string    `json:"district"`
	State             string    `json:"state"`
	Pincode           string    `json:"pincode"`
	Gst_Number        string    `json:"gst_number"`
	Brand             string    `json:"brand"`
	Model_Number      string    `json:"model_number"`
	Serial_Number     string    `json:"serial_number"`
	Issue_Description string    `json:"issue_description"`
	Due_Date          time.Time `json:"due_date"`
	Created_At        time.Time `json:"createdAt"`
	Updated_At        time.Time `json:"updatedAt"`
}

func (t *Ticket) UnmarshalJSON(data []byte) error {
	type Alias Ticket
	aux := &struct {
		Due_Date string `json:"due_date"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse the Due_Date string in RFC3339 format with timezone
	parsedDate, err := time.Parse(time.RFC3339, aux.Due_Date)
	if err != nil {
		return fmt.Errorf("invalid date format for Due_Date: %v", err)
	}
	t.Due_Date = parsedDate
	return nil
}
