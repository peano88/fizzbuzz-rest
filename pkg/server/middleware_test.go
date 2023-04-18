package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/peano88/fizzbazz-rest/pkg/model"
	"github.com/peano88/fizzbazz-rest/pkg/server/mocks"
	"github.com/peano88/fizzbazz-rest/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestValidationMiddleware_Ok(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "http://example.com?int1=2&int2=3&limit=7&str1=f&str2=b", nil)
	resp := httptest.NewRecorder()

	fbs := FizzBuzzServer{}

	handler := fbs.ValidationMiddleware(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		input := utils.FizzBuzzInputFromContext(r.Context())
		assert.NotNil(t, input)
	}))

	handler.ServeHTTP(resp, req)
}

func TestValidationMiddleware_Ko(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "http://example.com?int1=two&int2=3&limit=7&str1=f&str2=b", nil)
	resp := httptest.NewRecorder()

	fbs := FizzBuzzServer{}

	handler := fbs.ValidationMiddleware(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Fail(t, "should not call next handler")
	}))

	handler.ServeHTTP(resp, req)

	var appError model.ApplicationError
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&appError))
	assert.Equal(t, AppErrorTypeInput, appError.Type)
}

func TestStatisticsMiddleware_Ok(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	resp := httptest.NewRecorder()

	stats := mocks.NewFizzBuzzStats(t)
	stats.On("Increment", mock.AnythingOfType("*context.valueCtx"), 2, 3, 6, "f", "b").Return(nil)

	fbs := FizzBuzzServer{
		Stats: stats,
	}

	input := model.FizzBuzzInput{
		FizzBuzzInputStats: model.FizzBuzzInputStats{
			Int1:  2,
			Int2:  3,
			Limit: 6,
			Str1:  "f",
			Str2:  "b",
		},
	}

	handler := fbs.ToStatisticsMiddleware(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	}))

	handler.ServeHTTP(resp, req.WithContext(context.WithValue(req.Context(), model.InputKey, input)))
	stats.AssertExpectations(t)
}

func TestStatisticsMiddleware_Ko(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	resp := httptest.NewRecorder()

	stats := mocks.NewFizzBuzzStats(t)
	stats.On("Increment", mock.AnythingOfType("*context.valueCtx"), 2, 3, 6, "f", "b").Return(errors.New("dummy"))

	fbs := FizzBuzzServer{
		Stats: stats,
	}

	input := model.FizzBuzzInput{
		FizzBuzzInputStats: model.FizzBuzzInputStats{
			Int1:  2,
			Int2:  3,
			Limit: 6,
			Str1:  "f",
			Str2:  "b",
		},
	}

	toBeChanged := new(bool)
	*toBeChanged = false
	// In case of error in the middleware, next should be called anyway
	handler := fbs.ToStatisticsMiddleware(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		*toBeChanged = true
	}))

	handler.ServeHTTP(resp, req.WithContext(context.WithValue(req.Context(), model.InputKey, input)))
	assert.True(t, *toBeChanged)
	stats.AssertExpectations(t)
}
