package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ashermp9/fiskaly-test-task/config"
	"github.com/ashermp9/fiskaly-test-task/internal/app"
	"github.com/ashermp9/fiskaly-test-task/internal/ports"
	"github.com/ashermp9/fiskaly-test-task/internal/storage"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {

	var configFile string
	flag.StringVar(&configFile, "config", "config/local/config.yaml", "Path to the config file")
	flag.Parse()

	cfg := &config.Config{}
	err := config.LoadConfig(configFile, cfg)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        // Human-readable time format
	zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Capitalized and colored level
	zapConfig.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder      // Short caller format

	logger, err := zapConfig.Build()
	if err != nil {
		panic(fmt.Errorf("failed to initialize logger: %w", err))
	}
	defer func() { _ = logger.Sync() }()
	sugar := logger.Sugar()
	// Initialize components
	stor := storage.NewStorage()
	appService := app.NewAPIService(stor)

	// Set up and start the HTTP server
	server := ports.NewServer(sugar, appService, cfg.ServerAddress)
	go func() {
		sugar.Infof("Starting server on port %d", cfg.ServerAddress)
		if err := server.Run(); err != nil && err != http.ErrServerClosed {
			sugar.Fatalf("Failed to listen and serve: %v", err)
		}
	}()

	// Graceful shutdown
	gracefulShutdown(server, sugar)
}

func gracefulShutdown(server *ports.Server, logger *zap.SugaredLogger) {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	<-stopChan // Wait for interrupt signal

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Info("Shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("Server shutdown failed: %+v", err)
	} else {
		logger.Info("Server exited properly")
	}
}
