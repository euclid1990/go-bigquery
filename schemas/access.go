package schemas

type Access struct {
	Id       int    `json:"id" csv:"id" bigquery:"id"`
	UserId   int    `json:"user_id" csv:"user_id" bigquery:"user_id"`
	AccessAt string `json:"access_at" csv:"access_at" bigquery:"access_at"`
}
