package agent

import (
	"log/slog"
	"time"

	models "github.com/YagodinArtem/metric_collector.git/internal/model"
	"github.com/YagodinArtem/metric_collector.git/internal/service"
)

type Poller interface {
	PollWithInterval()
}

type DefaultPoller struct {
	pollInterval time.Duration
	mss          service.MetricStorageService[models.Metrics, string]
}

func NewPoller(
	pollIntervalSec int,
	mss service.MetricStorageService[models.Metrics, string]) Poller {

	return &DefaultPoller{
		pollInterval: time.Duration(pollIntervalSec) * time.Second,
		mss:          mss,
	}

}

func (p DefaultPoller) PollWithInterval() {

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

func (p DefaultPoller) poll() {
	// TODO: implement runtime metrics polling
}
