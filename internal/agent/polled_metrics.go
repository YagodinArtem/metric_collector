package agent

import (
	"math/rand"
	"reflect"
	"runtime"

	models "github.com/YagodinArtem/metric_collector.git/internal/model"
)

type PolledMetric struct {
	metricType string
	metricName string
	source     string
}

var PolledMetricList = []PolledMetric{
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
