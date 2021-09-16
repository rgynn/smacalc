package symbol

import (
	"fmt"
	"math"
	"testing"
)

func TestSymbolCalculate(t *testing.T) {

	type testcase struct {
		Calculator *SMACalculator
		Expected   float64
	}

	testcases := []testcase{
		{
			Calculator: &SMACalculator{
				Count:    1,
				Messages: []Message{{Price: 1}},
			},
			Expected: 0.5,
		},
		{
			Calculator: &SMACalculator{
				Count:    2,
				Messages: []Message{{Price: 1}, {Price: 2}},
			},
			Expected: 0.75,
		},
		{
			Calculator: &SMACalculator{
				Count:    3,
				Messages: []Message{{Price: 1}, {Price: 2}, {Price: 3}},
			},
			Expected: 1,
		},
	}

	results := make(chan Result, 1)

	for i, tc := range testcases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			tc.Calculator.Calculate(results)
			res := <-results
			if diff := math.Abs(res.Mean - tc.Expected); diff > 0.01 {
				t.Fatalf("expected result: %f, got: %f", tc.Expected, res.Mean)
			}
		})
	}
}
