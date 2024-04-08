package agent

import (
	"flag"
	"os"
	"strconv"
)

var endpointAddr string
var pollInterval int
var reportInterval int

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags() {
	// регистрируем переменную endpointAddr
	// как аргумент -a со значением localhost:8080 по умолчанию
	flag.StringVar(&endpointAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&pollInterval, "p", 2, "Poll interval in seconds")
	flag.IntVar(&reportInterval, "r", 10, "Report interval in seconds")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		endpointAddr = envRunAddr
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		v, err := strconv.Atoi(envReportInterval)
		if err != nil {
			panic(err)
		} else {
			reportInterval = v

		}
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		v, err := strconv.Atoi(envPollInterval)
		if err != nil {
			panic(err)
		} else {
			pollInterval = v
		}
	}
}
