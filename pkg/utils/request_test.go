package utils

import (
	"context"
	"reflect"
	"testing"

	"github.com/peano88/fizzbuzz-rest/pkg/model"
	"github.com/stretchr/testify/assert"
)

var emptyOutput = model.FizzBuzzStatisticsOutput{}

func TestStatisticsOutputFromString(t *testing.T) {

	buildExpected := func(n, m, limit int, fizz, buzz string, hits int64) model.FizzBuzzStatisticsOutput {
		return model.FizzBuzzStatisticsOutput{
			Parameters: model.FizzBuzzInputStats{
				Int1:  n,
				Int2:  m,
				Limit: limit,
				Str1:  fizz,
				Str2:  buzz,
			},
			Hits: hits,
		}
	}

	tests := []struct {
		label    string
		input    string
		hits     int64
		expected model.FizzBuzzStatisticsOutput
		wantErr  bool
	}{
		{"all good", "2-3-7-fzz-bzz", 9, buildExpected(2, 3, 7, "fzz", "bzz", 9), false},
		{"more than 5", "2-3-7-fzz-b-zz", 9, emptyOutput, true},
		{"int1 error", "two-3-7-fzz-bzz", 9, emptyOutput, true},
		{"int2 error", "2-three-7-fzz-bzz", 9, emptyOutput, true},
		{"limit error", "2-3-seven-fzz-bzz", 9, emptyOutput, true},
	}

	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			obtained, err := FizzBuzzStatisticsOutputFromString(tt.input, tt.hits)
			assert.Equal(t, tt.wantErr, err != nil, tt.label)
			assert.True(t, reflect.DeepEqual(tt.expected, obtained), tt.label)
		})
	}

}

func TestGetInputFromContext(t *testing.T) {
	input := model.FizzBuzzInput{
		FizzBuzzInputStats: model.FizzBuzzInputStats{
			Int1:  1,
			Int2:  2,
			Limit: 3,
			Str1:  "f",
			Str2:  "b",
		},
		Start: 4,
	}

	ctx := context.WithValue(context.TODO(), model.InputKey, input)
	assert.True(t, reflect.DeepEqual(input, FizzBuzzInputFromContext(ctx)))
}
