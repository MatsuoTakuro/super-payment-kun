package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"super-payment-kun/internal/config"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func OpenDB(ctx context.Context, cfg *config.Config) (*sqlx.DB, func(), error) {

	db, err := sql.Open("mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true",
			cfg.DBUser, cfg.DBPassword,
			cfg.DBHost, cfg.DBPort,
			cfg.DBName,
		),
	)
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, func() { _ = db.Close() }, err
	}
	log.Printf("connected to db (%s, %s:%d)", cfg.DBName, cfg.DBHost, cfg.DBPort)

	xdb := sqlx.NewDb(db, "mysql")
	return xdb, func() { _ = db.Close() }, nil
}
