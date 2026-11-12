package service

import (
	"errors"
	"sync"

	models "github.com/YagodinArtem/metric_collector.git/internal/model"
	"github.com/YagodinArtem/metric_collector.git/internal/repository"
)

type MetricStorageService[TYPE any, ID any] interface {
	SaveNewMetric(metric *TYPE) (bool, error)

	FindByID(ID) (*TYPE, error)
}

type MemStorageMetricStorageService struct {
	repo repository.Repository[models.Metrics, string]
	lock sync.RWMutex
}

func NewMemStorageMetricStorageService(repo repository.Repository[models.Metrics, string]) *MemStorageMetricStorageService {
	return &MemStorageMetricStorageService{
		repo: repo,
	}
}

func (m *MemStorageMetricStorageService) SaveNewMetric(metric *models.Metrics) (bool, error) {

	m.lock.Lock()
	defer m.lock.Unlock()

	existMetric, _ := m.repo.FindByID(metric.ID)

	if existMetric == nil {
		return m.repo.Save(metric)
	}

	if existMetric.MType != metric.MType {
		return false, errors.New("metric exist, but type is differ, change metric ID")
	}

	switch existMetric.MType {
	case models.Counter:

		delta := *metric.Delta + *existMetric.Delta
		existMetric.Delta = &delta

	case models.Gauge:
		existMetric.Value = metric.Value
	}

	_, err := m.repo.Update(existMetric)

	if err != nil {
		return false, err
	}

	return true, nil

}

func (m *MemStorageMetricStorageService) FindByID(id string) (*models.Metrics, error) {
	return m.repo.FindByID(id)
}
