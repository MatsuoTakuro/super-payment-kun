package router

import (
	"context"
	"fmt"
	"net/http"
	"super-payment-kun/internal/config"
	"super-payment-kun/internal/handler"
	"super-payment-kun/internal/pkg"
	"super-payment-kun/internal/repository"
	"super-payment-kun/internal/router/middleware"

	"github.com/go-chi/chi"
)

func NewRouter(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	db, cleanup, err := repository.OpenDB(ctx, cfg)
	if err != nil {
		return nil, cleanup, fmt.Errorf("failed to open db: %w", err)
	}

	clocker := pkg.RealClocker{}

	jwter, err := pkg.NewJWTer(clocker, cfg)
	if err != nil {
		return nil, cleanup, fmt.Errorf("failed to create JWTer: %w", err)
	}

	r := chi.NewRouter()

	_, err = repository.New(db, clocker)
	if err != nil {
		return nil, cleanup, fmt.Errorf("failed to create repository: %w", err)
	}

	_ = pkg.GetValidator()

	loginHdlr := handler.NewTestLogin(jwter)
	r.Post("/api/testlogin", loginHdlr.ServeHTTP)

	r.Route("/api/invoices", func(r chi.Router) {
		r.Use(middleware.AuthJWT(jwter))
	})

	return r, cleanup, nil
}
