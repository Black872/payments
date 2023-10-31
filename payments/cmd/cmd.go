package cmd

import (
	"context"
	"net"
	"os"
	"os/signal"
	"payments/config"
	"payments/db"
	"payments/repository"
	"payments/server"
	"payments/server/handlers"
	"payments/service"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func RunApp() (err error) {
	connDB, err := db.NewPostgresDB(db.DatabaseURL())
	if err != nil {
		return err
	}
	defer connDB.Close()

	// Create the base context with system signal handling.
	baseCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize server with base context for incoming requests, handlers, services, and repository.
	listener, err := net.Listen("tcp", config.ServerPort())
	if err != nil {
		return err
	}
	defer listener.Close()
	appHandlers := handlers.NewHandlers(service.NewService(repository.NewRepository(connDB)))
	appServer := server.NewServer(baseCtx, listener, appHandlers)

	// Create an error group to implement graceful shutdown.
	eg, egCtx := errgroup.WithContext(baseCtx)
	// Run a goroutine with the server start.
	eg.Go(func() error {
		return appServer.Run()
	})
	// Run a goroutine that waits for the system's signals.
	eg.Go(func() error {
		<-egCtx.Done()
		return appServer.Close(context.Background())
	})

	return eg.Wait()
}
