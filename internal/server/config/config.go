package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var RunAddr string
var LogLevel string
var StoreInterval int
var FileStoragePath string
var Restore bool
var DBConnect string

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func ParseFlags() error {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением: 8080 по умолчанию
	flag.StringVar(&RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&LogLevel, "l", "info", "log level")
	flag.IntVar(&StoreInterval, "i", 300, "store to file interval in seconds")
	flag.StringVar(&FileStoragePath, "f", "/tmp/metrics-db.json", "path to storage file")
	flag.BoolVar(&Restore, "r", true, "restore metrics from file")
	flag.StringVar(&DBConnect, "d", "host=localhost user=developer password=dev12", "DB connection string")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		RunAddr = envRunAddr
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		LogLevel = envLogLevel
	}

	var err error
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		StoreInterval, err = strconv.Atoi(envStoreInterval)
		if err != nil {
			return fmt.Errorf("invalid STORE_INTERVAL: %s", err)
		}
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		FileStoragePath = envFileStoragePath
	}

	if envRestore, ok := os.LookupEnv("RESTORE"); ok {
		if envRestore == "" || envRestore == "false" || envRestore == "0" {
			Restore = false
		} else {
			Restore = true
		}
	}

	if envDBConnect := os.Getenv("DATABASE_DSN"); envDBConnect != "" {
		DBConnect = envDBConnect
	}

	return nil
}
