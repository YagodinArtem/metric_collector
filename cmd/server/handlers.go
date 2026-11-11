package main

import (
	"log/slog"
	"net/http"

	models "github.com/YagodinArtem/metric_collector.git/internal/model"
)

const (
	headerContentTypeKey    = "Content-Type"
	contentTypeTextPlainVal = "text/plain"
)

func MetricRequestHandler(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodPost {
		slog.Error("unhandled request method", "method", request.Method)
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	contentType := request.Header.Get(headerContentTypeKey)

	if contentType != contentTypeTextPlainVal {
		slog.Error("unhandled content type", "content type", contentType)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	metricName := request.PathValue(metricNameKey)

	if metricName == "" {
		slog.Error("metric name not given")
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	metricType := request.PathValue(metricTypeKey)
	metricVal := request.PathValue(metricValueKey)

	if !models.MetricTypeIsCorrect(metricType) {
		slog.Info("metric type is incorrect", "given type", metricType)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if !models.MetricValIsSupported(metricVal, metricType) {
		slog.Info("metric val is incorrect", "given val", metricVal)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	metric := models.Metrics{
		ID:    metricName,
		MType: metricType,
		Delta: nil,
		Value: nil,
		Hash:  "",
	}

	metric.SetVal(metricVal)

	//TODO: save in storage
	slog.Info("Saved!", "metric", metric)

}
