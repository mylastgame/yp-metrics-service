package config

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	EndpointAddr   string
	PollInterval   int
	ReportInterval int
}

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func Create() (Config, error) {
	cfg := Config{
		EndpointAddr:   "",
		PollInterval:   0,
		ReportInterval: 0,
	}

	flag.StringVar(&cfg.EndpointAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&cfg.PollInterval, "p", 2, "Poll interval in seconds")
	flag.IntVar(&cfg.ReportInterval, "r", 10, "Report interval in seconds")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.EndpointAddr = envRunAddr
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		v, err := strconv.Atoi(envReportInterval)
		if err != nil {
			return Config{}, err
		} else {
			cfg.ReportInterval = v

		}
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		v, err := strconv.Atoi(envPollInterval)
		if err != nil {
			return Config{}, err
		} else {
			cfg.PollInterval = v
		}
	}

	return cfg, nil
}
