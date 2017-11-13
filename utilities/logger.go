package utilities

import (
	"fmt"
	"github.com/euclid1990/go-bigquery/configs"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func Log(level string, message string) {
	now := time.Now().Format(configs.LOG_FORMAT_DATE)

	logfile, err := os.OpenFile(configs.LOG_PATH+strings.ToLower(level)+now+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	defer logfile.Close()
	if err != nil {
		log.Fatalf("%s Open file log failed!", configs.LOG_CRITICAL)
	}
	prefix := fmt.Sprintf("%s: ", level)
	logger := log.New(io.MultiWriter(logfile, os.Stdout), prefix, log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println(message)
	if level == configs.LOG_CRITICAL || level == configs.LOG_ERROR {
		os.Exit(1)
	}
}

func Logf(level string, format string, value ...interface{}) {
	Log(level, fmt.Sprintf(format, value...))
}
