package input

import "testing"
import "fmt"

type TestCase struct {
	input   Input
	success func(interface{}) bool
}

func TestEval(t *testing.T) {
	inputs := []TestCase{

		// Value
		TestCase{
			input: Input{
				Type:  "integer",
				Value: 1,
			},
			success: func(v interface{}) bool {
				return v == 1
			},
		},

		TestCase{
			input: Input{
				Type:  "double",
				Value: 1.5,
			},
			success: func(v interface{}) bool {
				return v == 1.5
			},
		},

		// Or
		TestCase{
			input: Input{
				Type: "string",
				Or:   []interface{}{"aaa", "bbb"},
			},
			success: func(v interface{}) bool {
				return v == "aaa" || v == "bbb"
			},
		},

		// Range
		TestCase{
			input: Input{
				Type: "integer",
				Max:  float64p(1.0),
				Min:  float64p(-1.0),
			},
			success: func(v interface{}) bool {
				num := v.(int)
				return num >= -1 && num <= 1
			},
		},

		TestCase{
			input: Input{
				Type: "double",
				Max:  float64p(1.0),
				Min:  float64p(-1.0),
			},
			success: func(v interface{}) bool {
				num := v.(float64)
				return num >= -1 && num <= 1
			},
		},
	}

	for i, in := range inputs {
		if v := in.input.Eval(); !in.success(v) {
			t.Errorf("test case %d failed. input is %v, evaluation value is %v.", i, in.input, v)
		}
	}

}

func float64p(value interface{}) *float64 {
	var v float64
	var ok bool
	if v, ok = value.(float64); !ok {
		panic(fmt.Sprintf("fail to convert float64: %v", value))
	}

	return &v
}
