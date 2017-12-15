package schemas

import (
	"cloud.google.com/go/civil"
	"time"
)

type DateTime struct {
	time.Time
}

// Convert the internal date as CSV string
func (date *DateTime) MarshalCSV() (string, error) {
	return date.Time.Format("2006-01-02 15:04:05"), nil
}

// Convert the CSV string as internal date
func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse("2006-01-02 15:04:05", csv)
	if err != nil {
		return err
	}
	return nil
}

// Convert Datetime struct to string
func (date DateTime) ToString() string {
	return date.Time.Format("2006-01-02 15:04:05")
}

// Convert Datetime struct to civil.Datetime
func (date DateTime) ToCivil() civil.DateTime {
	return civil.DateTimeOf(date.Time)
}
