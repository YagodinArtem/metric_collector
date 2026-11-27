package handler

import (
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"
	"sort"

	models "github.com/YagodinArtem/metric_collector.git/internal/model"
	"github.com/YagodinArtem/metric_collector.git/internal/service"
	"github.com/go-chi/chi/v5"
)

const (
	HeaderContentTypeKey    = "Content-Type"
	ContentTypeTextPlainVal = "text/plain"
	ContentTypeTextHTMLVal  = "text/html; charset=utf-8"

	htmlTemplate = `<html>
<head>
    <title>Metrics list</title>
</head>
<body>
    <h1>Metrics list</h1>
    <div>
        {{range .}}
            <pre>{{.}}</pre>
        {{else}}
            <p>no metrics found</p>
        {{end}}
    </div>
</body>
</html>`
)

var tmpl = template.Must(template.New("metrics").Parse(htmlTemplate))

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

	metricName := chi.URLParam(request, models.MetricNameKey)

	if metricName == "" {
		slog.Error("metric name not given")
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	metricType := chi.URLParam(request, models.MetricTypeKey)
	metricVal := chi.URLParam(request, models.MetricValueKey)

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

	updatedMetric, err := h.mss.SaveNewMetric(&metric)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte("unable to save data"))
		return
	}

	if updatedMetric != nil {
		slog.Info("saved!", "metric", updatedMetric)
		writer.WriteHeader(http.StatusOK)
	} else {
		slog.Error("metric not saved", "metric", metric)
		writer.WriteHeader(http.StatusInternalServerError)
	}

}

func (h *MetricHandler) AllMetricsListHTML(writer http.ResponseWriter, _ *http.Request) {

	allMetrics, _ := h.mss.FindAll()

	sort.Slice(allMetrics, func(i, j int) bool {
		return allMetrics[i].ID < allMetrics[j].ID
	})

	jsonList := make([]string, 0, len(allMetrics))

	for _, metric := range allMetrics {

		marshal, err := json.Marshal(metric)

		if err != nil {
			slog.Error("unable to json marshall metric", "error", err)
			continue
		}

		jsonList = append(jsonList, string(marshal))

	}

	err := tmpl.Execute(writer, jsonList)

	writer.Header().Set(HeaderContentTypeKey, ContentTypeTextHTMLVal)

	if err != nil {
		slog.Error("error write html page", "error ", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (h *MetricHandler) GetSingleMetric(writer http.ResponseWriter, request *http.Request) {

	id := chi.URLParam(request, models.MetricNameKey)

	metric, _ := h.mss.FindByID(id)

	if metric == nil {
		http.NotFound(writer, request)
		return
	}

	metricVal, _ := models.MetricValToStr(metric)

	writer.Header().Set(HeaderContentTypeKey, ContentTypeTextPlainVal)
	_, _ = writer.Write([]byte(metricVal))

}
