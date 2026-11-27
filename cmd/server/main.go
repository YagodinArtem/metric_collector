package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/YagodinArtem/metric_collector.git/internal"
	"github.com/YagodinArtem/metric_collector.git/internal/config"
)

func main() {
	err := RunServer()
	if err != nil {
		panic(err)
	}
}

func RunServer() error {

	cfg, configLoadErr := config.LoadServerConfig()

	if cfg == nil || configLoadErr != nil {
		return fmt.Errorf("server config loading error: %s", configLoadErr)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	router := internal.BuildRouter(internal.BuildMemStorageService())

	slog.Info("server is running")
	err := http.ListenAndServe(cfg.Address, router)

	if err != nil {
		slog.Info("error caught when running server", "error", err)
		return err
	}

	return nil
}
