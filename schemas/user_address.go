package schemas

import (
	"errors"
	"fmt"
	"strings"
)

type UserAddress struct {
	Status  string `json:"current" csv:"current" bigquery:"current"`
	City    string `json:"city" csv:"city" bigquery:"city"`
	Country string `json:"country" csv:"country" bigquery:"country"`
}

// Convert the a nested and repeated field UserAddress as CSV string
func (adr *UserAddress) MarshalCSV() (string, error) {
	return fmt.Sprintf("%s, %s, %s", adr.Status, adr.City, adr.Country), nil
}

// Convert the CSV string as a nested and repeated field UserAddress
func (adr *UserAddress) UnmarshalCSV(csv string) (err error) {
	s := strings.Split(csv, ", ")
	if len(s) != 3 {
		return errors.New("The format of the address is not correct")
	}
	adr.Status, adr.City, adr.Country = s[0], s[1], s[2]
	return nil
}
