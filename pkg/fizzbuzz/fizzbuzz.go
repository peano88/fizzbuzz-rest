package fizzbuzz

import (
	"strconv"

	"github.com/peano88/fizzbuzz-rest/pkg/model"
)

// Fizzbuzz is a function creating a sequence based on the input
// namely, the sequence starts with input.start and stop at input.limit
// each element can be replaced by input.str1 and/or input.str2 if the number
// is a multiple of int1 and/or int2
func Fizzbuzz(input model.FizzBuzzInput) []string {
	result := []string{}
	for i := input.Start; i <= input.Limit; i++ {
		output := ""
		if i%input.Int1 == 0 {
			output += input.Str1
		}
		if i%input.Int2 == 0 {
			output += input.Str2
		}
		if output == "" {
			output = strconv.Itoa(i)
		}
		result = append(result, output)
	}
	return result
}
