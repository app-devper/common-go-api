package form

type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
}
