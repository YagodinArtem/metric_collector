package handler_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/YagodinArtem/metric_collector.git/internal"
	"github.com/YagodinArtem/metric_collector.git/internal/handler"
	"github.com/stretchr/testify/require"
)

func TestMetricHandler_MetricRequestHandler(t *testing.T) {
	service := internal.BuildMemStorageService()
	mux := internal.BuildMux(service)

	tests := []struct {
		testName string
		pattern  string
		expected int
		header   bool
	}{
		{
			"save gauge success test",
			"/update/gauge/test_gauge/64.0",
			200,
			true,
		},

		{
			"save counter success test",
			"/update/counter/test_counter/256",
			200,
			true,
		},

		{
			"bad metric name failed test",
			"/update/counter/256",
			404,
			true,
		},

		{
			"gauge bad value failed test",
			"/update/gauge/test_gauge_bad_value/hello_bug",
			400,
			true,
		},

		{
			"counter bad value failed test",
			"/update/counter/counter_bad_val/hello_bugs",
			400,
			true,
		},

		{
			"counter bad value failed test",
			"/update/counter/counter_bad_val/hello_bugs",
			400,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodPost, tt.pattern, nil)

			writer := httptest.NewRecorder()

			if tt.header {
				request.Header.Set(handler.HeaderContentTypeKey, handler.ContentTypeTextPlainVal)
			}

			mux.ServeHTTP(writer, request)

			result := writer.Result()

			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					fmt.Print(err)
				}
			}(result.Body)

			require.Equal(t, tt.expected, result.StatusCode)

		})
	}
}

func TestMetricHandler_Saving(t *testing.T) {
	service := internal.BuildMemStorageService()
	mux := internal.BuildMux(service)

	requests := []string{
		"/update/gauge/test_gauge/64.0",
		"/update/gauge/test_gauge/128.0",
		"/update/gauge/test_gauge/256.0",

		"/update/gauge/test_gauge2/256.0",
		"/update/gauge/test_gauge2/512.0",

		"/update/counter/test_counter/32",
		"/update/counter/test_counter/64",
		"/update/counter/test_counter/32",

		"/update/counter/test_counter2/32",
		"/update/counter/test_counter2/32",
	}

	for _, url := range requests {
		req := httptest.NewRequest(http.MethodPost, url, nil)
		req.Header.Set(handler.HeaderContentTypeKey, handler.ContentTypeTextPlainVal)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
	}

	gauge, _ := service.FindByID("test_gauge")
	gauge2, _ := service.FindByID("test_gauge2")
	counter, _ := service.FindByID("test_counter")
	counter2, _ := service.FindByID("test_counter2")

	require.Equal(t, 256.0, *gauge.Value)
	require.Equal(t, 512.0, *gauge2.Value)

	require.Equal(t, int64(128), *counter.Delta)
	require.Equal(t, int64(64), *counter2.Delta)

}
