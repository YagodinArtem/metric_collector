package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/YagodinArtem/metric_collector.git/internal"
	models "github.com/YagodinArtem/metric_collector.git/internal/model"
)

func main() {
	err := RunServer()
	if err != nil {
		panic(err)
	}
}

func RunServer() error {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)
	mux := internal.BuildMux(internal.BuildMemStorageService())

	slog.Info("server is running ")
	err := http.ListenAndServe(models.ServerPath, mux)

	if err != nil {
		slog.Info("error caught when running server", "error", err)
		return err
	}

	return nil
}
