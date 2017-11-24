package utilities

import (
	"encoding/json"
	"github.com/euclid1990/go-bigquery/configs"
	"github.com/gocarina/gocsv"
	"github.com/joho/godotenv"
	"os"
	"time"
)

func LoadEnv(file string) {
	if file == "" {
		file = ".env"
	}
	err := godotenv.Load(file)
	if err != nil {
		Logf(configs.LOG_CRITICAL, "Error loading %v file", file)
	}
}

func Random(min, max int) int {
	rndSrc := CreateRndSrc()
	return rndSrc.Intn(max-min) + min
}

func RandomDateTimeBetween(afterObj time.Time, beforeObj time.Time) time.Time {
	afterUnix := afterObj.Unix()
	beforeUnix := beforeObj.Unix()

	diff := beforeUnix - afterUnix

	rndSrc := CreateRndSrc()
	sec := rndSrc.Int63n(diff) + afterUnix

	return time.Unix(sec, 0)
}

func WriteCsv(file string, data interface{}) (bool, error) {
	csvFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()
	if err := gocsv.MarshalFile(data, csvFile); err != nil {
		return false, err
	}
	return true, nil
}

func WriteJson(file string, data interface{}) (bool, error) {
	jsonFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	buff, errJson := json.MarshalIndent(data, "", "	")
	if errJson != nil {
		return false, errJson
	}
	if _, err := jsonFile.Write(buff); err != nil {
		panic(err)
	}
	return true, nil
}
