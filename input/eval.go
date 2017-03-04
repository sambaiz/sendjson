package input

import (
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Input struct {
	Type       string        `json:"type"`  // string, boolean, integer, double, or time
	Value      interface{}   `json:"value"` // fixed value
	Or         []interface{} `json:"or"`    // randomly selected
	Max        *float64      `json:"max"`   // range, ignored except for integer and double
	Min        *float64      `json:"min"`
	TimeFormat string        `json:"time_format"` // time string or unix_epoch, ignored except for time
}

// eval input then return value
func (in Input) Eval() interface{} {

	evals := []func() (interface{}, bool){
		in.evalValue,
		in.evalOr,
		in.evalRange,
		in.evalTimeFormat,
	}

	for _, ev := range evals {
		if v, ok := ev(); ok {
			return v
		}
	}

	return nil
}

func typeCheck(typeName string, value interface{}) bool {
	switch typeName {
	case "string":
		if _, ok := value.(string); ok {
			return true
		}
	case "boolean":
		if _, ok := value.(bool); ok {
			return true
		}
	case "integer":
		if _, ok := value.(int); ok {
			return true
		}
	case "double":
		if _, ok := value.(float64); ok {
			return true
		}
	}

	return false
}

func (in Input) evalValue() (interface{}, bool) {

	if in.Value == nil {
		return nil, false
	}

	if ok := typeCheck(in.Type, in.Value); !ok {
		return nil, false
	}

	return in.Value, true
}

func (in Input) evalOr() (interface{}, bool) {

	if in.Or == nil {
		return nil, false
	}

	if v := in.Or[rand.Intn(len(in.Or))]; v == nil {
		return nil, false
	} else if ok := typeCheck(in.Type, v); !ok {
		return nil, false
	} else {
		return v, true
	}
}

func (in Input) evalRange() (interface{}, bool) {

	if in.Max == nil || in.Min == nil {
		return nil, false
	}

	max := *(in.Max)
	min := *(in.Min)

	v := (max-min)*rand.Float64() + min

	switch in.Type {
	case "integer":
		return int(math.Floor(v + 0.5)), true
	case "double":
		return v, true
	}

	return nil, false
}

func (in Input) evalTimeFormat() (interface{}, bool) {

	if in.TimeFormat == "" {
		return nil, false
	}

	switch in.TimeFormat {
	case "unix_epoch":
		return now().Unix(), true
	default:
		return now().Format(in.TimeFormat), true
	}
}

var now = func() time.Time {
	return time.Now()
}
