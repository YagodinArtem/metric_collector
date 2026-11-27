package internal

import (
	"github.com/YagodinArtem/metric_collector.git/internal/handler"
	"github.com/YagodinArtem/metric_collector.git/internal/model"
	"github.com/YagodinArtem/metric_collector.git/internal/repository"
	"github.com/YagodinArtem/metric_collector.git/internal/service"
	"github.com/go-chi/chi/v5"
)

func BuildMemStorageService() service.MetricStorageService[models.Metrics, string] {
	return service.NewMemStorageMetricStorageService(BuildMemStorage())
}

func BuildMemStorage() repository.Repository[models.Metrics, string] {
	return repository.NewMemStorage()
}

func BuildRouter(mss service.MetricStorageService[models.Metrics, string]) chi.Router {
	router := chi.NewRouter()

	postMetricPattern := "/update/{type}/{name}/{value}"
	getMetricPattern := "/value/{type}/{name}"
	defaultPath := "/"

	metricHandler := handler.NewMetricHandler(mss)

	router.Post(postMetricPattern, metricHandler.MetricRequestHandler)
	router.Get(getMetricPattern, metricHandler.GetSingleMetric)
	router.Get(defaultPath, metricHandler.AllMetricsListHTML)

	return router
}
