package user

type User struct {
	ID       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" db:"password" validate:"required,min=6"`
}
