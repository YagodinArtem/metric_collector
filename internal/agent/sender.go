package agent

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/YagodinArtem/metric_collector.git/internal/handler"
	models "github.com/YagodinArtem/metric_collector.git/internal/model"
	"github.com/YagodinArtem/metric_collector.git/internal/service"
)

const metricUpdateURLTemplate = "http://%s/update/%s/%s/%s"

type Sender interface {
	SendWithInterval()
}

type DefaultSender struct {
	interval time.Duration
	address  string
	mss      service.MetricStorageService[models.Metrics, string]
	client   http.Client
}

func NewSender(sendIntervalSec int,
	address string,
	mss service.MetricStorageService[models.Metrics, string]) Sender {

	client := http.Client{
		Timeout: time.Second * 5,
	}

	return &DefaultSender{
		interval: time.Duration(sendIntervalSec) * time.Second,
		address:  address,
		mss:      mss,
		client:   client,
	}
}

func (d DefaultSender) SendWithInterval() {

	for {
		time.Sleep(d.interval)
		d.send()
	}

}

func (d DefaultSender) send() {
	metrics, _ := d.mss.FindAll()

	for _, m := range metrics {

		valStr, err := models.MetricValToStr(&m)

		if err != nil {
			slog.Error("metric val is nil", "metric", m)
			continue
		}

		url := fmt.Sprintf(metricUpdateURLTemplate, d.address, m.MType, m.ID, valStr)

		post, err := d.client.Post(url, handler.ContentTypeTextPlainVal, nil)

		if err != nil {
			slog.Error("metric POST error", "error", err)
		} else {
			slog.Info("metric POST result", "metric", m, "response code", post.StatusCode)
			_ = post.Body.Close()
			d.postHandle(&m)
		}

	}

}

func (d DefaultSender) postHandle(m *models.Metrics) {

	switch m.ID {
	case "PollCount":
		_, _ = d.mss.Delete(m)
	}
}
