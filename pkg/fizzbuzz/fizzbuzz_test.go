package fizzbuzz

import (
	"testing"

	"github.com/peano88/fizzbazz-rest/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestFizzBuzz(t *testing.T) {

	buildInput := func(n, m, limit int, fizz, buzz string, start int) model.FizzBuzzInput {
		return model.FizzBuzzInput{
			FizzBuzzInputStats: model.FizzBuzzInputStats{
				Int1:  n,
				Int2:  m,
				Limit: limit,
				Str1:  fizz,
				Str2:  buzz,
			},
			Start: start,
		}
	}

	tests := []struct {
		label    string
		input    model.FizzBuzzInput
		expected []string
	}{
		{"normal case", buildInput(2, 3, 7, "fizz", "buzz", 1), []string{"1", "fizz", "buzz", "fizz", "5", "fizzbuzz", "7"}},
		{"not reachable fizz", buildInput(7, 3, 6, "fizz", "buzz", 2), []string{"2", "buzz", "4", "5", "buzz"}},
		{"not reachable buzz", buildInput(3, 7, 6, "fuzz", "buzz", 1), []string{"1", "2", "fuzz", "4", "5", "fuzz"}},
		{"no fizz", buildInput(6, 7, 5, "fizz", "buzz", 1), []string{"1", "2", "3", "4", "5"}},
		{"empty", buildInput(2, 3, 0, "fizz", "buzz", 1), []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			output := Fizzbuzz(model.FizzBuzzInput(tt.input))
			assert.Equal(t, tt.expected, output)
		})
	}

}
