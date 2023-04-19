package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/httplog"
	"github.com/peano88/fizzbuzz-rest/pkg/utils"
	"github.com/peano88/fizzbuzz-rest/pkg/validation"
)

// ValidationMiddleware is an HTTP middleware which runs a set of validation on the 
// query parameter of the request. It forwards a modified context.Context to the next handler
// obtained by inserting the validated set of input parameters
func (fbs *FizzBuzzServer) ValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		forwardContext, err := validation.NewFizzBuzzValidator().RunValidations(r)
		if err != nil {
			oplog := httplog.LogEntry(r.Context())
			oplog.Err(fmt.Errorf("validation error: %w", err)).Msg("")
			validationApplicationError(rw, r, err)
			return
		}

		next.ServeHTTP(rw, r.WithContext(forwardContext))
	})
}

// ToStatisticsMiddleware is an HTTP middleware sending the set of input parameters to the statistics component.
// the set is retrieved via the request context.Context. If an error arises, then the error is logged, but the 
// next handler is called anyway
func (fbs *FizzBuzzServer) ToStatisticsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		input := utils.FizzBuzzInputFromContext(r.Context())

		if err := fbs.Stats.Increment(r.Context(), input.Int1, input.Int2, input.Limit, input.Str1, input.Str2); err != nil {
			oplog := httplog.LogEntry(r.Context())
			oplog.Err(fmt.Errorf("error incrementing stats: %w", err)).Msg("")
		}

		next.ServeHTTP(rw, r)
	})
}
