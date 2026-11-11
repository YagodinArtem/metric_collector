package models

import (
	"slices"
	"strconv"
)

const (
	Counter = "counter"
	Gauge   = "gauge"
)

var validMetricTypes = []string{Counter, Gauge}

// NOTE: Не усложняем пример, вводя иерархическую вложенность структур.
// Органичиваясь плоской моделью.
// Delta и Value объявлены через указатели,
// что бы отличать значение "0", от не заданного значения
// и соответственно не кодировать в структуру.

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Hash  string   `json:"hash,omitempty"`
}

func (m *Metrics) SetVal(givenVal string) {

	switch m.MType {
	case Counter:
		counterVal, _ := strconv.ParseInt(givenVal, 10, 64)
		m.Delta = &counterVal
	case Gauge:
		gaugeVal, _ := strconv.ParseFloat(givenVal, 64)
		m.Value = &gaugeVal
	}
}

func MetricTypeIsCorrect(givenType string) bool {
	return slices.Contains(validMetricTypes, givenType)
}

func MetricValIsSupported(givenVal string, givenType string) bool {
	switch givenType {
	case Counter:
		_, err := strconv.ParseInt(givenVal, 10, 64)
		return err == nil
	case Gauge:
		_, err := strconv.ParseFloat(givenVal, 64)
		return err == nil
	default:
		return false
	}
}
