package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"super-payment-kun/internal/config"
	"super-payment-kun/internal/router"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatalf("failed to terminate server normally: %v", err)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.HttpPort))
	if err != nil {
		return fmt.Errorf("failed to listen port %d: %w", cfg.HttpPort, err)
	}

	r, cleanup, err := router.NewRouter(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to create router: %w", err)
	}
	defer cleanup()

	s := NewServer(l, r)
	return s.Run(ctx)

}
