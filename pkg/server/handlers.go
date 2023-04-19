package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/httplog"
	"github.com/peano88/fizzbuzz-rest/pkg/fizzbuzz"
	"github.com/peano88/fizzbuzz-rest/pkg/model"
	"github.com/peano88/fizzbuzz-rest/pkg/utils"
)

const (
    // Header key for content type
	ContentTypeHeader = "Content-Type"
    // Header value for JSON content type
	JSONContentType   = "application/json"

    // Limit of items for a single response to the fizzbuzz sequence
    // if [start,limit] is an interval containing more than paginationLimit elements,
    // the response will contains a sequence of paginationLimit elements and
    // a link to the request completing/extending the sequence
	paginationLimit = 65536
)

// GetFizzBuzzHandler is the handler for the /fizzbuzz endpoint under method GET. 
// It expects int1, int2, limit, str1, str2 query parameters and allows the optional
// start parameter. The response is a fizz-buzz-alike sequence from start to limit (start
// is defaulted to 1 if not provided). If the sequence consists of more than paginationLimit
// elements, the response is paginated. 
func (fbs *FizzBuzzServer) GetFizzBuzzHandler(rw http.ResponseWriter, r *http.Request) {
	input := utils.FizzBuzzInputFromContext(r.Context())
	oplog := httplog.LogEntry(r.Context())

	output := model.FizzBuzzOutput{}
	pagination := (input.Limit - input.Start) > paginationLimit

	if pagination {
		output.Next = fmt.Sprintf("/fizzbuzz?int1=%d&int2=%d&limit=%d&start=%d&str1=%s&str2=%s", input.Int1, input.Int2, input.Limit, input.Start+paginationLimit, input.Str1, input.Str2)
		input.Limit = input.Start + paginationLimit - 1
	}

	output.Sequence = fizzbuzz.Fizzbuzz(input)

	respPayload, err := json.Marshal(&output)
	if err != nil {
		oplog.Err(fmt.Errorf("error marshaling response: %w", err)).Msg("")
		jsonApplicationError(rw, r)
		return
	}
	rw.Header().Add(ContentTypeHeader, JSONContentType)
	rw.Write(respPayload)
	rw.WriteHeader(http.StatusOK)
}

// GetStatisticsHandler is the handler for the /statistics endpoint under method GET. The
// response is the set of input parameters most requested. If two sets share the same request count,
// then the set returned is the first by reversed lexicographical order. Please note that the start parameter
// of GET /fizzbuzz is not considered in the input parameter set; furthermore, only a validated set (i.e. a set 
// where the input parameters are complaint with the validations) is considered for the statistics 
func (fbs *FizzBuzzServer) GetStatisticsHandler(rw http.ResponseWriter, r *http.Request) {
	res, err := fbs.Stats.Stats(r.Context())
	if err != nil {
		statisticsApplicationError(rw, r, err)
		return
	}

	respPayload, err := json.Marshal(&res)
	if err != nil {
		oplog := httplog.LogEntry(r.Context())
		oplog.Err(fmt.Errorf("error marshaling response: %w", err)).Msg("")
		jsonApplicationError(rw, r)
		return
	}

	rw.Header().Add(ContentTypeHeader, JSONContentType)
	rw.WriteHeader(http.StatusOK)
	rw.Write(respPayload)
}
