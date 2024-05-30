package model

import (
	"time"
)

type VendorBankAccount struct {
	ID            string     `db:"id"` // uuid
	Vendor        Vendor     `db:"vendor"`
	BankName      string     `db:"bank_name"`
	BranchName    string     `db:"branch_name"`
	AccountNumber string     `db:"account_number"`
	AccountName   string     `db:"account_name"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at"`
}
