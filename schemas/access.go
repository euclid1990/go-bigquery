package schemas

import (
	"cloud.google.com/go/civil"
)

type Access struct {
	Id       int            `json:"id" csv:"id" bigquery:"id"`
	UserId   int            `json:"user_id" csv:"user_id" bigquery:"user_id"`
	AccessAt civil.DateTime `json:"access_at" csv:"access_at" bigquery:"access_at"`
}
