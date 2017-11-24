package schemas

type Access struct {
	Id         int      `csv:"id"`
	UserId     int      `csv:"user_id"`
	AccessedAt DateTime `csv:"accessed_at"`
}
