# Golang and Google BigQuery

Quickstart using Go-lang with Google Cloud Platform to select/insert/update data in BigQuery

## Config project

- Create `.env` file

```bash
$ cp .env.example .env
```

- Create credential (*Service account key*) at `https://console.cloud.google.com/apis/credentials?project=[Project_ID]` and put in inside root of `go-bigquery` with following name

```
google_application_credentials.json
```


## Install dependencies

```bash
$ glide install

```

## Start program

```bash
$ go run main.go --exec [action]
```

## Reference

- https://cloud.google.com/bigquery/docs/reference/libraries
- https://cloud.google.com/bigquery/create-simple-app-api#bigquery-simple-app-build-service-go
