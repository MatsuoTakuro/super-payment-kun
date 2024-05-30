package repository

import (
	"fmt"
	"super-payment-kun/internal/pkg"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	db    *sqlx.DB
	clock pkg.Clocker
}

func New(db *sqlx.DB, c pkg.Clocker) (*repository, error) {

	if db == nil || c == nil {
		return nil, fmt.Errorf("failed to create repository")
	}
	return &repository{
		db:    db,
		clock: c,
	}, nil
}
