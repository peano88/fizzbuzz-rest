package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/peano88/fizzbuzz-rest/pkg/model"
	"github.com/peano88/fizzbuzz-rest/pkg/statistics"
	"github.com/peano88/fizzbuzz-rest/pkg/validation"
)

const (
    // ApplicationError type for json-marshaling error
	AppErrorTypeJSON    = "/fizzbuzz/errors/json"
    // ApplicationError type for parsing-related error
	AppErrorTypeParsing = "/fizzbuzz/errors/parse"
    // ApplicationError type for statistics availability error
	AppErrorTypeStats   = "/fizzbuzz/errors/stats"
    // ApplicationError type for error generated during validation of the input parameters
	AppErrorTypeInput   = "/fizzbuzz/errors/input"
)

func jsonApplicationError(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusInternalServerError)
	appError := model.ApplicationError{
		Type:     AppErrorTypeJSON,
		Title:    "error marshaling response",
		Status:   strconv.Itoa(http.StatusInternalServerError),
		Instance: middleware.GetReqID(r.Context()),
	}

	appErrorPayload, err := json.Marshal(&appError)
	if err != nil {
		oplog := httplog.LogEntry(r.Context())
		oplog.Err(fmt.Errorf("app error marshaling issue: %w", err)).Msg("")
		return
	}

	rw.Header().Add(ContentTypeHeader, JSONContentType)
	rw.Write(appErrorPayload)
}

func validationApplicationError(rw http.ResponseWriter, r *http.Request, err error) {
	rw.WriteHeader(http.StatusBadRequest)

	var valErr validation.ValidationError
	if !errors.As(err, &valErr) {
		oplog := httplog.LogEntry(r.Context())
		oplog.Err(fmt.Errorf("validation error expected")).Msg("")
		return
	}

	appError := model.ApplicationError{
		Type:     AppErrorTypeInput,
		Title:    valErr.Error(),
		Status:   strconv.Itoa(http.StatusBadRequest),
		Detail:   valErr.Constraint(),
		Instance: middleware.GetReqID(r.Context()),
	}

	appErrorPayload, err := json.Marshal(&appError)
	if err != nil {
		oplog := httplog.LogEntry(r.Context())
		oplog.Err(fmt.Errorf("application error marshaling issue: %w", err)).Msg("")
		return
	}

	rw.Header().Add(ContentTypeHeader, JSONContentType)
	rw.Write(appErrorPayload)
}

func statisticsApplicationError(rw http.ResponseWriter, r *http.Request, err error) {
	appError := model.ApplicationError{
		Type:     AppErrorTypeStats,
		Title:    "internal issue retrieving statistics",
		Instance: middleware.GetReqID(r.Context()),
	}
	if errors.Is(err, statistics.NoStatsAvailable{}) {
		rw.WriteHeader(http.StatusServiceUnavailable)
		appError.Status = strconv.Itoa(http.StatusServiceUnavailable)
		appError.Detail = "no previous request available"

	} else {
		rw.WriteHeader(http.StatusInternalServerError)
		appError.Status = strconv.Itoa(http.StatusInternalServerError)
	}
	oplog := httplog.LogEntry(r.Context())
	oplog.Err(fmt.Errorf("error getting statistics: %w", err)).Msg("")

	appErrorPayload, err := json.Marshal(&appError)
	if err != nil {
		oplog := httplog.LogEntry(r.Context())
		oplog.Err(fmt.Errorf("application error marshaling issue: %w", err)).Msg("")
		return
	}

	rw.Header().Add(ContentTypeHeader, JSONContentType)
	rw.Write(appErrorPayload)
}
