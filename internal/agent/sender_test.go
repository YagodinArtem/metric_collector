package agent

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	models "github.com/YagodinArtem/metric_collector.git/internal/model"
	"github.com/YagodinArtem/metric_collector.git/internal/repository"
	"github.com/YagodinArtem/metric_collector.git/internal/service"
	"github.com/stretchr/testify/require"
)

func TestDefaultSender_SendWithInterval(t *testing.T) {

	storage := repository.NewMemStorage()
	mss := service.NewMemStorageMetricStorageService(storage)

	mux := http.NewServeMux()

	postMetricPattern := fmt.Sprintf("POST /update/{%s}/{%s}/{%s}", models.MetricTypeKey, models.MetricNameKey, models.MetricValueKey)

	var receiveCount atomic.Int64
	mux.HandleFunc(postMetricPattern, func(writer http.ResponseWriter, request *http.Request) {
		receiveCount.Add(1)
		writer.WriteHeader(http.StatusOK)
	})

	ts := httptest.NewServer(mux)

	serverAddr := strings.TrimPrefix(ts.URL, "http://")

	poller := NewPoller(1, mss)
	sender := NewSender(2, serverAddr, mss)

	go poller.SinglePollWithInterval()
	go sender.SendWithInterval()

	time.Sleep(time.Duration(3) * time.Second)

	require.Greater(t, receiveCount.Load(), int64(5))

}
