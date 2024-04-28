package counter

import "strconv"

const Key = "counter"

type ValType int64

type Counter struct {
	Title string
	Val   ValType
}

func New(title string, val ValType) *Counter {
	return &Counter{Title: title, Val: val}
}

func FromString(title, value string) (*Counter, error) {
	v, err := strconv.ParseInt(value, 0, 64)
	if err != nil {
		return &Counter{}, err
	}

	return New(title, ValType(v)), nil
}
