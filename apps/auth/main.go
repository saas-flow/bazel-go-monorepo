package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func init() {
	// Load file .env jika ada
	if err := godotenv.Load(); err != nil {
		zap.L().Info("No .env file found, skipping...")
	}
}

func main() {
	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// Handle OS signals untuk graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nReceived shutdown signal, exiting...")
		cancel()
	}()

	// Start Uber FX application
	app := fx.New(Module)

	// Run application
	if err := app.Start(ctx); err != nil {
		zap.L().Fatal("Failed to start application", zap.Error(err))
	}

	// Wait for context cancellation
	<-app.Done()

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := app.Stop(shutdownCtx); err != nil {
		zap.L().Fatal("Failed to stop app:", zap.Error(err))
	}
}
