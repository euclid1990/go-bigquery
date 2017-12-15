# Golang and Google BigQuery

Quickstart using Go-lang with Google Cloud Platform to select/insert/update data in BigQuery

![Go-lang and Bigquery web service](web/assets/images/demo.gif)

## Config project

- Create `.env` file

```bash
$ cp .env.example .env
```

- Open `.env` file and put your GCP Project Id: `GCP_PROJECT_ID`

- Create credential (*Service account key*) at `https://console.cloud.google.com/apis/credentials?project=[Project_ID]` and put in inside root of `go-bigquery` with following name

```
google_application_credentials.json
```

- Please enable billing for this BigQuery Project. If you not connect to active billing, you will falling back to error:
```
Error 403: Billing has not been enabled for this project. Enable billing at https://console.cloud.google.com/billing.
```


## Install dependencies

```bash
$ glide install

```

## Start program

```bash
$ go run main.go --exec [action]
```

## Start web application

```bash
$ go run main.go -exec web
```

Go to: http://localhost:8000/


## Reference

- https://godoc.org/cloud.google.com/go/bigquery
- https://cloud.google.com/bigquery/docs/reference/libraries
- https://cloud.google.com/bigquery/create-simple-app-api#bigquery-simple-app-build-service-go
