package agent

import (
	"log/slog"
	"runtime"
	"time"

	models "github.com/YagodinArtem/metric_collector.git/internal/model"
	"github.com/YagodinArtem/metric_collector.git/internal/service"
)

const (
	RandomSource  = "random"
	RuntimeSource = "runtime"
)

type MetricGetter[INPUT any, OUTPUT any] interface {
	GetMetricValue(INPUT, PolledMetric) OUTPUT
}

type Poller interface {
	InfinityPoll()
	SinglePollWithInterval()
}

var runtimeMetricGetter = RuntimeMetricGetter{}
var randomMetricGetter = RandomMetricGetter{}

type DefaultMetricsPoller struct {
	pollInterval time.Duration
	mss          service.MetricStorageService[models.Metrics, string]
	polledList   []PolledMetric
}

func NewPoller(
	pollIntervalSec int,
	mss service.MetricStorageService[models.Metrics, string]) Poller {

	return &DefaultMetricsPoller{
		pollInterval: time.Duration(pollIntervalSec) * time.Second,
		mss:          mss,
		polledList:   PolledMetricList,
	}

}

func (p DefaultMetricsPoller) InfinityPoll() {

	for {
		p.SinglePollWithInterval()
	}

}

func (p DefaultMetricsPoller) SinglePollWithInterval() {
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

func (p DefaultMetricsPoller) poll() {

	for _, polledMetric := range p.polledList {

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
