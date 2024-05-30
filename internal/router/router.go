package router

import (
	"context"
	"fmt"
	"net/http"
	"super-payment-kun/internal/config"
	"super-payment-kun/internal/pkg"
	"super-payment-kun/internal/repository"

	"github.com/go-chi/chi"
)

func NewRouter(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	db, cleanup, err := repository.OpenDB(ctx, cfg)
	if err != nil {
		return nil, cleanup, fmt.Errorf("failed to open db: %w", err)
	}

	clocker := pkg.RealClocker{}

	r := chi.NewRouter()

	_, err = repository.New(db, clocker)
	if err != nil {
		return nil, cleanup, fmt.Errorf("failed to create repository: %w", err)
	}

	_ = pkg.GetValidator()

	return r, cleanup, nil
}
