package schemas

import (
	"cloud.google.com/go/civil"
)

type User struct {
	Id        int            `json:"id" csv:"id" bigquery:"id"`
	Name      string         `json:"name" csv:"name" bigquery:"name"`
	Age       int            `json:"age" csv:"age" bigquery:"age"`
	Email     string         `json:"email" csv:"email" bigquery:"email"`
	Gender    string         `json:"gender" csv:"gender" bigquery:"gender"`
	Address   UserAddress    `json:"address" csv:"address" bigquery:"address"`
	CreatedAt civil.DateTime `json:"created_at" csv:"created_at" bigquery:"created_at"`
}
