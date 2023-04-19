package utils

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/peano88/fizzbuzz-rest/pkg/model"
)

// FizzBuzzInputFromContext returns the model.FizzBuzzInput contained in the context.Context.
// Will PANIC if not present
func FizzBuzzInputFromContext(ctx context.Context) model.FizzBuzzInput {
	return ctx.Value(model.InputKey).(model.FizzBuzzInput)
}

// FizzBuzzStatisticsOutputFromString splits the s string using the model.Separator and creates the model.FizzBuzzStatisticsOutput 
// using the separated tokens
func FizzBuzzStatisticsOutputFromString(s string, hits int64) (model.FizzBuzzStatisticsOutput, error) {
	tokens := strings.Split(s, model.Separator)

	if len(tokens) != 5 {
		return model.FizzBuzzStatisticsOutput{}, fmt.Errorf("input parameters string is incorrect: %s", s)
	}
	int1, err := strconv.Atoi(tokens[0])
	if err != nil {
		return model.FizzBuzzStatisticsOutput{}, fmt.Errorf("int1 can't be parsed: %w", err)
	}
	int2, err := strconv.Atoi(tokens[1])
	if err != nil {
		return model.FizzBuzzStatisticsOutput{}, fmt.Errorf("int2 can't be parsed: %w", err)
	}
	limit, err := strconv.Atoi(tokens[2])
	if err != nil {
		return model.FizzBuzzStatisticsOutput{}, fmt.Errorf("limit can't be parsed: %w", err)
	}

	return model.FizzBuzzStatisticsOutput{
		Parameters: model.FizzBuzzInputStats{
			Int1:  int1,
			Int2:  int2,
			Limit: limit,
			Str1:  tokens[3],
			Str2:  tokens[4],
		},
		Hits: hits,
	}, nil
}
