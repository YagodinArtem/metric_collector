package agent

import (
	"log/slog"
	"math/rand"
	"reflect"
	"runtime"
	"time"

	models "github.com/YagodinArtem/metric_collector.git/internal/model"
	"github.com/YagodinArtem/metric_collector.git/internal/service"
)

type Poller interface {
	PollWithInterval()
}

const (
	RandomSource  = "random"
	RuntimeSource = "runtime"
)

type MetricGetter[INPUT any, OUTPUT any] interface {
	GetMetricValue(INPUT, PolledMetric) OUTPUT
}

type RuntimeMetricGetter struct {
}

type RandomMetricGetter struct {
}

func (r RuntimeMetricGetter) GetMetricValue(m *runtime.MemStats, pm PolledMetric) float64 {
	v := reflect.ValueOf(m).Elem().FieldByName(pm.metricName)
	if !v.IsValid() {
		return 0
	}

	switch v.Kind() {
	case reflect.Uint64, reflect.Uint32:
		return float64(v.Uint())
	case reflect.Float64:
		return v.Float()
	default:
		return 0
	}
}

func (r RandomMetricGetter) GetMetricValue(_ *any, _ PolledMetric) float64 {
	return rand.Float64()
}

var runtimeMetricGetter = RuntimeMetricGetter{}
var randomMetricGetter = RandomMetricGetter{}

type PolledMetric struct {
	metricType string
	metricName string
	source     string
}

var metricList = []PolledMetric{
	{models.Gauge, "BuckHashSys", RuntimeSource},
	{models.Gauge, "Frees", RuntimeSource},
	{models.Gauge, "GCCPUFraction", RuntimeSource},
	{models.Gauge, "HeapAlloc", RuntimeSource},
	{models.Gauge, "HeapIdle", RuntimeSource},
	{models.Gauge, "HeapInuse", RuntimeSource},
	{models.Gauge, "HeapObjects", RuntimeSource},
	{models.Gauge, "HeapReleased", RuntimeSource},
	{models.Gauge, "HeapSys", RuntimeSource},
	{models.Gauge, "LastGC", RuntimeSource},
	{models.Gauge, "Lookups", RuntimeSource},
	{models.Gauge, "MCacheInuse", RuntimeSource},
	{models.Gauge, "MCacheSys", RuntimeSource},
	{models.Gauge, "MSpanInuse", RuntimeSource},
	{models.Gauge, "MSpanSys", RuntimeSource},
	{models.Gauge, "Mallocs", RuntimeSource},
	{models.Gauge, "NextGC", RuntimeSource},
	{models.Gauge, "NumForcedGC", RuntimeSource},
	{models.Gauge, "NumGC", RuntimeSource},
	{models.Gauge, "OtherSys", RuntimeSource},
	{models.Gauge, "PauseTotalNs", RuntimeSource},
	{models.Gauge, "StackInuse", RuntimeSource},
	{models.Gauge, "StackSys", RuntimeSource},
	{models.Gauge, "Sys", RuntimeSource},
	{models.Gauge, "TotalAlloc", RuntimeSource},

	{models.Gauge, "RandomValue", RandomSource},
}

type DefaultMetricsPoller struct {
	pollInterval time.Duration
	mss          service.MetricStorageService[models.Metrics, string]
}

func NewPoller(
	pollIntervalSec int,
	mss service.MetricStorageService[models.Metrics, string]) Poller {

	return &DefaultMetricsPoller{
		pollInterval: time.Duration(pollIntervalSec) * time.Second,
		mss:          mss,
	}

}

func (p DefaultMetricsPoller) PollWithInterval() {

	for {
		p.poll()

		poll := &models.Metrics{
			ID:    "PollCount",
			MType: "counter",
			Delta: new(int64(1)),
			Value: nil,
			Hash:  "",
		}

		update, _ := p.mss.SaveNewMetric(poll)

		if update != nil {
			curPoll, _ := p.mss.FindByID(poll.ID)
			slog.Info("current poll iter: ", "poll", *curPoll.Delta)
		}

		time.Sleep(p.pollInterval)
	}

}

func (p DefaultMetricsPoller) poll() {

	for _, polledMetric := range metricList {

		slog.Info("now poll", "polled metric name", polledMetric.metricName)

		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		switch polledMetric.source {

		case RuntimeSource:

			if polledMetric.metricType == models.Gauge {

				metricVal := runtimeMetricGetter.GetMetricValue(&m, polledMetric)

				metricToSave := &models.Metrics{
					ID:    polledMetric.metricName,
					MType: polledMetric.metricType,
					Delta: nil,
					Value: &metricVal,
					Hash:  "",
				}

				_, err := p.mss.SaveNewMetric(metricToSave)

				if err != nil {
					slog.Error("error during saving new metric", "error", err)
				}
			}

		case RandomSource:

			if polledMetric.metricType == models.Gauge {
				metricVal := randomMetricGetter.GetMetricValue(nil, polledMetric)

				metricToSave := &models.Metrics{
					ID:    polledMetric.metricName,
					MType: polledMetric.metricType,
					Delta: nil,
					Value: &metricVal,
					Hash:  "",
				}

				_, err := p.mss.SaveNewMetric(metricToSave)

				if err != nil {
					slog.Error("error during saving new metric", "error", err)
				}
			}

		}
	}
}
