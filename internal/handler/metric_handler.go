package handler

import (
	"log/slog"
	"net/http"

	models "github.com/YagodinArtem/metric_collector.git/internal/model"
	"github.com/YagodinArtem/metric_collector.git/internal/service"
)

const (
	HeaderContentTypeKey    = "Content-Type"
	ContentTypeTextPlainVal = "text/plain"
)

type MetricHandler struct {
	mss service.MetricStorageService[models.Metrics, string]
}

func NewMetricHandler(mss service.MetricStorageService[models.Metrics, string]) *MetricHandler {
	return &MetricHandler{
		mss: mss,
	}
}

func (h *MetricHandler) MetricRequestHandler(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodPost {
		slog.Error("unhandled request method", "method", request.Method)
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	contentType := request.Header.Get(HeaderContentTypeKey)

	if contentType != ContentTypeTextPlainVal {
		slog.Error("unhandled content type", "content type", contentType)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	metricName := request.PathValue(models.MetricNameKey)

	if metricName == "" {
		slog.Error("metric name not given")
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	metricType := request.PathValue(models.MetricTypeKey)
	metricVal := request.PathValue(models.MetricValueKey)

	if !models.MetricTypeIsCorrect(metricType) {
		slog.Error("metric type is incorrect", "given type", metricType)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if !models.MetricValIsSupported(metricVal, metricType) {
		slog.Error("metric val is incorrect", "given val", metricVal)
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

	success, err := h.mss.SaveNewMetric(&metric)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte("unable to save data"))
		return
	}

	if success {
		slog.Info("saved!", "metric", metric)
		writer.WriteHeader(http.StatusOK)
	} else {
		slog.Error("metric not saved", "metric", metric)
		writer.WriteHeader(http.StatusInternalServerError)
	}

}
