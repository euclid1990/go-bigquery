package utilities

import (
	"errors"
	"fmt"
	"github.com/dmgk/faker"
	"github.com/euclid1990/go-bigquery/configs"
	s "github.com/euclid1990/go-bigquery/schemas"
	"github.com/icrowley/fake"
	"os"
	"strconv"
	"time"
)

const (
	PERIOD_USER_CREATED_AT  = 13
	FROM_ACCESS_MONTH_AGO   = 12
	GENERATE_ACCESS_TIMEOUT = 100
)

var (
	UserAccessLog = make(map[int][]int64)
	CountryLog    = make(map[string]bool)
	TOTAL_USER    = 0
	TOTAL_ACCESS  = 0
	TOTAL_COUNTRY = 20
)

func InitTotalRecord() {
	TOTAL_USER, _ = strconv.Atoi(os.Getenv("TOTAL_USER_RECORD"))
	TOTAL_ACCESS, _ = strconv.Atoi(os.Getenv("TOTAL_ACCESS_RECORD"))
}

func GenrateDummyData(filetype string) int64 {
	InitTotalRecord()
	// Generate user data
	userData := GenerateUser(TOTAL_USER)
	userFile := fmt.Sprintf(configs.DATA_FORMAT_FILE_NAME, configs.DATA_INPUT_PATH+configs.DATA_INPUT_USER, filetype)
	// Generate user access log data
	accessData := GenerateAccess(TOTAL_ACCESS)
	accessFile := fmt.Sprintf(configs.DATA_FORMAT_FILE_NAME, configs.DATA_INPUT_PATH+configs.DATA_INPUT_ACCESS, filetype)
	switch filetype {
	case configs.DATA_TYPE_CSV:
		WriteCsv(userFile, userData)
		WriteCsv(accessFile, accessData)
	case configs.DATA_TYPE_JSON:
		WriteJson(userFile, userData)
		WriteJson(accessFile, accessData)
	}
	return int64(TOTAL_USER + TOTAL_ACCESS)
}

func GenerateUser(total int) []*s.User {
	user := make([]*s.User, total)
	now := time.Now()
	// Subtract PERIOD_USER_CREATED_AT months to current datetime
	startCreatedAt := now.AddDate(0, -1*PERIOD_USER_CREATED_AT, 0)
	endCreatedAt := now.AddDate(0, 0, -1)
	uniqueCountries := make([]string, 0)
	for i := 0; i < total; i++ {
		country := fake.Country()
		if len(uniqueCountries) < TOTAL_COUNTRY {
			_, ok := CountryLog[country]
			if !ok {
				uniqueCountries = append(uniqueCountries, country)
			}
			CountryLog[country] = true
		} else {
			n := Random(0, TOTAL_COUNTRY-1)
			country = uniqueCountries[n]
		}
		user[i] = &s.User{
			Id:        i,
			Name:      faker.Name().Name(),
			Age:       Random(15, 50),
			Email:     faker.Internet().Email(),
			Gender:    fake.Gender(),
			Address:   s.UserAddress{Status: "current", City: faker.Address().City(), Country: country},
			CreatedAt: s.DateTime{RandomDateTimeBetween(startCreatedAt, endCreatedAt)}.ToCivil(),
		}
	}
	return user
}

func GenerateAccess(total int) []*s.Access {
	access := make([]*s.Access, total)
	now := time.Now()
	// Subtract 2 Month to current datetime
	startCreatedAt := now.AddDate(0, -1*FROM_ACCESS_MONTH_AGO, 0)
	for i := 0; i < total; i++ {
		userId := Random(0, TOTAL_USER)
		accessAt := createUniqueAccessByUser(&userId, startCreatedAt, now)
		access[i] = &s.Access{
			Id:       i,
			UserId:   userId,
			AccessAt: accessAt.ToCivil(),
		}
	}
	return access
}

func createUniqueAccessArray(userId *int) {
	_, ok := UserAccessLog[*userId]
	if !ok {
		UserAccessLog[*userId] = make([]int64, 0)
	}
}

func createUniqueAccessByUser(userId *int, start, end time.Time) s.DateTime {
	createUniqueAccessArray(userId)
	totalRandomAccessAt := int(end.Sub(start) / time.Second)
	startAt := time.Now()
	for {
		if len(UserAccessLog[*userId]) == totalRandomAccessAt {
			*userId = Random(0, TOTAL_USER)
			createUniqueAccessArray(userId)
		}
		accessDatetime := s.DateTime{RandomDateTimeBetween(start, end)}
		accessDatetimeStamp := accessDatetime.Unix()
		unique := true
		for _, aDt := range UserAccessLog[*userId] {
			if aDt == accessDatetimeStamp {
				unique = false
			}
		}
		if unique {
			UserAccessLog[*userId] = append(UserAccessLog[*userId], accessDatetimeStamp)
			return accessDatetime
		}
		if time.Now().Sub(startAt) > (GENERATE_ACCESS_TIMEOUT * time.Second) {
			panic(errors.New("Timeout exceeded while dump new accessed at time"))
		}
	}
}
