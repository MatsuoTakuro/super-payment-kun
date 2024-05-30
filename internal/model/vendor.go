package model

import "time"

type Vendor struct {
	ID                 string     `db:"id"` // uuid
	Company            Company    `db:"company"`
	Name               string     `db:"name"`
	RepresentativeName string     `db:"representative_name"`
	PhoneNumber        string     `db:"phone_number"`
	ZipCode            string     `db:"zip_code"`
	Address            string     `db:"address"`
	CreatedAt          time.Time  `db:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at"`
	DeletedAt          *time.Time `db:"deleted_at"`
}
