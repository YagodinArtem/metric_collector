package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/YagodinArtem/metric_collector.git/internal/agent"
	"github.com/YagodinArtem/metric_collector.git/internal/config"
	"github.com/YagodinArtem/metric_collector.git/internal/repository"
	"github.com/YagodinArtem/metric_collector.git/internal/service"
)

func main() {

	err := runAgent()

	if err != nil {
		panic(err)
	}

	slog.Info("agent is stopped")

}

func runAgent() error {

	cfg, configLoadErr := config.AgentFlagConfig()

	if cfg == nil || configLoadErr != nil {
		return fmt.Errorf("agent config loading error: %s", configLoadErr)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	storage := repository.NewMemStorage()
	mss := service.NewMemStorageMetricStorageService(storage)
	poller := agent.NewPoller(cfg.PollInterval, mss)
	sender := agent.NewSender(cfg.ReportInterval, cfg.Address, mss)

	slog.Info("agent is running")

	go poller.InfinityPoll()
	sender.SendWithInterval()

	return nil
}
