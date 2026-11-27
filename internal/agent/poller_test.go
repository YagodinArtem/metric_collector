package agent

import (
	"testing"
	"time"

	models "github.com/YagodinArtem/metric_collector.git/internal/model"
	"github.com/YagodinArtem/metric_collector.git/internal/repository"
	"github.com/YagodinArtem/metric_collector.git/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultMetricsPoller_PollWithInterval(t *testing.T) {

	storage := repository.NewMemStorage()
	mss := service.NewMemStorageMetricStorageService(storage)
	poller := NewPoller(3, mss)

	cur := time.Now().UnixMilli()

	poller.SinglePollWithInterval()

	elapsed := time.Now().UnixMilli() - cur

	assert.True(t, elapsed > 3*1000)

	for _, metric := range PolledMetricList {

		localSavedMetric, err := mss.FindByID(metric.metricName)

		assert.Nil(t, err)
		assert.NotNil(t, localSavedMetric)

		if metric.metricType == models.Gauge {
			require.NotNil(t, localSavedMetric.Value)
			assert.GreaterOrEqual(t, *localSavedMetric.Value, float64(0))
		}

	}

}
