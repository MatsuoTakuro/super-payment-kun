package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	srv *http.Server
	l   net.Listener
}

func NewServer(l net.Listener, router http.Handler) *Server {
	return &Server{
		srv: &http.Server{Handler: router},
		l:   l,
	}
}

func (s *Server) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// Run the server and check the result of the shutdown
		log.Printf("server is getting started on %s", s.l.Addr())
		if err := s.srv.Serve(s.l); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("failed to close: %+v", err)
		}
		return nil
	})

	<-ctx.Done() // Wait for the interrupt signal or termination signal
	if err := s.srv.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: +%v", err)
	}

	return eg.Wait()
}
