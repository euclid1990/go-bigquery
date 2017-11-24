package schemas

type User struct {
	Id        int         `json:"id" csv:"id"`
	Name      string      `json:"name" csv:"name"`
	Age       int         `json:"age" csv:"age"`
	Email     string      `json:"email" csv:"email"`
	Gender    string      `json:"gender" csv:"gender"`
	Address   UserAddress `json:"address" csv:"address"`
	CreatedAt DateTime    `json:"created_at" csv:"created_at"`
}
