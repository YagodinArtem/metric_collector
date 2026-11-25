package internal

import (
	"fmt"
	"net/http"

	"github.com/YagodinArtem/metric_collector.git/internal/handler"
	"github.com/YagodinArtem/metric_collector.git/internal/model"
	"github.com/YagodinArtem/metric_collector.git/internal/repository"
	"github.com/YagodinArtem/metric_collector.git/internal/service"
)

func BuildMemStorageService() service.MetricStorageService[models.Metrics, string] {
	return service.NewMemStorageMetricStorageService(BuildMemStorage())
}

func BuildMemStorage() repository.Repository[models.Metrics, string] {
	return repository.NewMemStorage()
}

func BuildMux(mss service.MetricStorageService[models.Metrics, string]) *http.ServeMux {
	mux := http.NewServeMux()
	postMetricPattern := fmt.Sprintf("POST /update/{%s}/{%s}/{%s}", models.MetricTypeKey, models.MetricNameKey, models.MetricValueKey)

	metricHandler := handler.NewMetricHandler(mss)

	mux.HandleFunc(postMetricPattern, metricHandler.MetricRequestHandler)

	return mux
}
