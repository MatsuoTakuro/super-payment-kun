package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

// TODO: use transaction management if needed
// runInTX executes the given function in a transaction.
func (r *repository) runInTX(ctx context.Context, f func(context.Context, *sql.Tx) error) error {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	err = f(ctx, tx)
	if err != nil {
		log.Println("db rollback...")
		if rbErr := tx.Rollback(); rbErr != nil {
			return errors.Join(err, fmt.Errorf("failed to rollback: %w", rbErr))
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}
