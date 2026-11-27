package config

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Address string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
}

type AgentConfig struct {
	Address        string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	PollInterval   int    `env:"POLL_INTERVAL" envDefault:"10"`
	ReportInterval int    `env:"REPORT_INTERVAL" envDefault:"2"`
}

func LoadServerConfig() (*ServerConfig, error) {

	_ = godotenv.Load(".env.local")

	var cfg ServerConfig
	err := env.Parse(&cfg)

	if err != nil {
		return nil, fmt.Errorf("server config loading error: %w", err)
	}

	return &cfg, nil
}

func ServerFlagConfig() (*ServerConfig, error) {

	serverAddress := flag.String("a", "localhost:8080", "server address")

	flag.Parse()

	return &ServerConfig{
		Address: *serverAddress,
	}, nil
}

func LoadAgentConfig() (*AgentConfig, error) {

	_ = godotenv.Load(".env.local")

	var cfg AgentConfig
	err := env.Parse(&cfg)

	if err != nil {
		return nil, fmt.Errorf("agent config loading error: %w", err)
	}

	return &cfg, nil
}

func AgentFlagConfig() (*AgentConfig, error) {

	agentAddress := flag.String("a", "localhost:8080", "server address")
	reportInterval := flag.String("r", "10", "report interval ")
	pollInterval := flag.String("p", "2", "poll interval ")

	report, _ := strconv.ParseInt(*reportInterval, 10, 64)
	poll, _ := strconv.ParseInt(*pollInterval, 10, 64)

	flag.Parse()

	return &AgentConfig{
		Address:        *agentAddress,
		PollInterval:   int(poll),
		ReportInterval: int(report),
	}, nil

}
