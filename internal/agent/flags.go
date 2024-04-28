package agent

import "flag"

var endpointAddr string
var pollInterval int64
var reportInterval int64

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags() {
	// регистрируем переменную endpointAddr
	// как аргумент -a со значением localhost:8080 по умолчанию
	flag.StringVar(&endpointAddr, "a", "localhost:8080", "address and port to run server")
	flag.Int64Var(&pollInterval, "p", 2, "Poll interval in seconds")
	flag.Int64Var(&reportInterval, "r", 10, "Report interval in seconds")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
}
