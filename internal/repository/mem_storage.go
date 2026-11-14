package repository

import (
	"errors"
	"sync"

	models "github.com/YagodinArtem/metric_collector.git/internal/model"
)

type MemStorage struct {
	lock    sync.RWMutex
	storage map[string]*models.Metrics
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		storage: make(map[string]*models.Metrics),
	}
}

func (s *MemStorage) Save(t *models.Metrics) (bool, error) {

	s.lock.Lock()
	defer s.lock.Unlock()

	if _, exists := s.storage[t.ID]; exists {
		return false, errors.New("metric already exists")
	}

	s.storage[t.ID] = t

	return true, nil
}

func (s *MemStorage) SaveAll(metrics []models.Metrics) error {

	for _, metric := range metrics {

		_, err := s.Save(&metric)

		if err != nil {
			return err
		}

	}

	return nil
}

func (s *MemStorage) FindByID(id string) (*models.Metrics, error) {
	s.lock.RLock()

	defer s.lock.RUnlock()

	metric, exist := s.storage[id]

	if exist {
		return metric, nil
	}

	return nil, nil

}

func (*MemStorage) FindAll() ([]models.Metrics, error) {
	panic("implement me")
}

func (s *MemStorage) Update(t *models.Metrics) (bool, error) {

	if t == nil {
		return false, errors.New("metric is nil")
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	if _, exists := s.storage[t.ID]; !exists {
		return false, errors.New("metric not found by id")
	}

	s.storage[t.ID] = t
	return true, nil

}

func (*MemStorage) DeleteByID(id string) (bool, error) {
	panic("implement me")
}

func (*MemStorage) DeleteAll() (bool, error) {
	panic("implement me")
}
