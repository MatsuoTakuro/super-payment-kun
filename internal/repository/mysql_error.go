package repository

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
)

const (
	// WARN: Uppercase and lowercase letters are not distinguished (due to the setting in utf8mb4_unicode_ci).
	// https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html
	MySQLDuplicateEntryErrCode       = 1062
	MySQLForeignKeyConstraintErrCode = 1452
)

var mysqlErr *mysql.MySQLError

func isDuplicateEntryErr(err error) bool {
	return errors.As(err, &mysqlErr) && mysqlErr.Number == MySQLDuplicateEntryErrCode
}

func isForeignKeyConstraintErr(err error) bool {
	return errors.As(err, &mysqlErr) && mysqlErr.Number == MySQLForeignKeyConstraintErrCode
}

var noRowErr error = sql.ErrNoRows
