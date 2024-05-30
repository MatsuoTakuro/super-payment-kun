package model

import "time"

type Role int8

const (
	Normal Role = 0
	Admin  Role = 1
)

type User struct {
	ID        string     `db:"id"` // uuid
	Company   Company    `db:"company"`
	Name      string     `db:"name"`
	Email     string     `db:"email"`
	Role      Role       `db:"role"`
	Password  string     `db:"password"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (u *User) IsAdmin() bool {
	return u.Role == Admin
}

func (u *User) CanCreateInvoice() bool {
	return u.IsAdmin()
}

func (u *User) CanReadAllInvoicesForCompany() bool {
	return u.IsAdmin()
}
