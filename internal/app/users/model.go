package loans

import "github.com/roguexray007/loan-app/internal/app/base"

const (
	TableUser  = "users"
	EntityUser = "user"
)

type User struct {
	base.Model
	Username string `json:"username"`
	Pass     string `json:"pass"`
}

// TableName sets the insert table name
func (u *User) TableName() string {
	return TableUser
}

// EntityName gives the display name of the entity
func (u *User) EntityName() string {
	return EntityUser
}
