package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

const (
	metricTypeKey  = "type"
	metricNameKey  = "name"
	metricValueKey = "value"
	serverPath     = "localhost:8080"
)

func RunServer() error {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)
	mux := buildMux()

	slog.Info("server is running ")
	err := http.ListenAndServe(serverPath, mux)

	if err != nil {
		slog.Info("error caught when running server", "error", err)
		return err
	}

	return nil
}

func buildMux() *http.ServeMux {
	mux := http.NewServeMux()
	pattern := fmt.Sprintf("POST /update/{%s}/{%s}/{%s}", metricTypeKey, metricNameKey, metricValueKey)
	mux.HandleFunc(pattern, MetricRequestHandler)

	return mux
}
