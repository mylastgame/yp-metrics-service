package gauge

import (
	"fmt"
	"strconv"
)

const Key = "gauge"

type ValType float64

type Gauge struct {
	Title string
	Val   ValType
}

func New(title string, val ValType) *Gauge {
	return &Gauge{Title: title, Val: val}
}

func FromString(title, value string) (*Gauge, error) {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return &Gauge{}, err
	}

	return New(title, ValType(v)), nil
}

func (g *Gauge) String() string {
	return fmt.Sprintf("%s: %f", g.Title, g.Val)
}
