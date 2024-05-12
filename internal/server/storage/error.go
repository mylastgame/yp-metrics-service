package storage

import (
	"errors"
	"fmt"
)

const (
	BadMetricType = 1 + iota
	KeyNotExists
	BadValue
)

type StorageError struct {
	Code int
	Msg  string
}

var ErrorNotExists = errors.New("value not exists")

func NewStorageError(code int, t string, k string) *StorageError {
	return &StorageError{Code: code, Msg: createMessage(code, t, k)}
}

func createMessage(code int, t string, k string) string {
	switch code {
	case BadMetricType:
		return fmt.Sprintf("wrong metric type: %s", t)
	case KeyNotExists:
		return fmt.Sprintf("key: %s not exists", k)
	case BadValue:
		return fmt.Sprintf("Bad value: %s for %s", k, t)
	default:
		return fmt.Sprintf("error: %d: %s, %s", code, t, k)
	}
}

func (e StorageError) Error() string {
	return e.Msg
}
