package model

// Separator is s string used for the concatenation of the fields
// of the input parameters
const Separator string = "-"

// InputContextKey is a specific type for a key of a value in a context.Context
type InputContextKey int

// InputKey is the key to use when adding a FizzBuzzInput in a context.Context
const InputKey InputContextKey = iota

// FizzBuzzInputStats is a subset of the fizzbuzz input parameters, used to store
// and calculate statistics afterwards
type FizzBuzzInputStats struct {
	// every multiple of int1 will be replaced by str1 or str1str2
	Int1 int
	// every multiple of int2 will be replaced by str2 or str1str2
	Int2 int
	// upper limit (inclusive) of the fizzbuzz sequence
	Limit int
	// the string used for replacing every multiple of int1
	Str1 string
	// the string used for replacing every multiple of int2
	Str2 string
}

// FizzBuzzInput collects all input parameters for the generation of a fizzbuzz sequence
// it embedes the structure FizzBuzzInputStats
type FizzBuzzInput struct {
	FizzBuzzInputStats
	// the number used to start the fizzbuzz sequence
	Start int
}

// FizzBuzzOutput is the structure returned by the the /fizzbuzz endpoint: the fizzbuzz Sequence
// and optionally a Next link, pointing to the next sequence, if the request underwent the pagination
type FizzBuzzOutput struct {
	// fizzBuzz sequence
	Sequence []string
	// next pagination request link
	Next string `json:"next,omitempty"`
}

// FizzBuzzStatisticsOutput is the structure returned by the /statistics endpoint: the most used
// input parameters set and the number of times that it has been requested
type FizzBuzzStatisticsOutput struct {
	// Set of Input parameters of /fizzbuzz endpoint, used to calculate the statistics
	Parameters FizzBuzzInputStats
	// Number of times the Parameters set has been requested
	Hits int64
}

// ApplicationError is the structure returned in case of error by the two endpoints
type ApplicationError struct {
	// URI formatted type of the error
	Type string `json:"err_type"`
	// Brief description, human-readable, error description
	Title string
	// HTTP status
	Status string
	// a more complete description of the error
	Detail string `json:"detail,omitempty"`
	// application identifier of the request generating the error
	Instance string
}
