package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/peano88/fizzbazz-rest/pkg/model"
	"github.com/peano88/fizzbazz-rest/pkg/server/mocks"
	"github.com/peano88/fizzbazz-rest/pkg/statistics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFizzBuzzHandler_OK(t *testing.T) {

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)

	input := model.FizzBuzzInput{
		FizzBuzzInputStats: model.FizzBuzzInputStats{
			Int1:  2,
			Int2:  3,
			Limit: 7,
			Str1:  "f",
			Str2:  "b",
		},
		Start: 1,
	}

	fbs := FizzBuzzServer{}

	fbs.GetFizzBuzzHandler(resp, req.WithContext(context.WithValue(req.Context(), model.InputKey, input)))

	var output model.FizzBuzzOutput
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&output))
	assert.Empty(t, output.Next)
	assert.True(t, reflect.DeepEqual([]string{"1", "f", "b", "f", "5", "fb", "7"}, output.Sequence), "%v", output.Sequence)
}

func TestGetFizzBuzzHandler_OK_WithPagination(t *testing.T) {

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)

	input := model.FizzBuzzInput{
		FizzBuzzInputStats: model.FizzBuzzInputStats{
			Int1:  2,
			Int2:  3,
			Limit: 65539,
			Str1:  "f",
			Str2:  "b",
		},
		Start: 1,
	}

	fbs := FizzBuzzServer{}

	fbs.GetFizzBuzzHandler(resp, req.WithContext(context.WithValue(req.Context(), model.InputKey, input)))

	var output model.FizzBuzzOutput
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&output))
	assert.NotEmpty(t, output.Next)
	assert.Equal(t, 65536, len(output.Sequence))
}

func TestGetStatistics_Ok(t *testing.T) {
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)

	stats := mocks.NewFizzBuzzStats(t)
	toReturn := model.FizzBuzzStatisticsOutput{
		Parameters: model.FizzBuzzInputStats{
			Int1:  2,
			Int2:  3,
			Limit: 7,
			Str1:  "f",
			Str2:  "b",
		},
		Hits: 9,
	}

	stats.On("Stats", req.Context()).Return(toReturn, nil)

	fbs := FizzBuzzServer{
		Stats: stats,
	}

	fbs.GetStatisticsHandler(resp, req)
	var output model.FizzBuzzStatisticsOutput
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&output))
	assert.Equal(t, toReturn, output)
}

func TestGetStatistics_Ko_noResults(t *testing.T) {
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)

	stats := mocks.NewFizzBuzzStats(t)

	stats.On("Stats", req.Context()).Return(model.FizzBuzzStatisticsOutput{}, statistics.NoStatsAvailable{})

	fbs := FizzBuzzServer{
		Stats: stats,
	}

	fbs.GetStatisticsHandler(resp, req)
	assert.Equal(t, http.StatusServiceUnavailable, resp.Result().StatusCode)
	var output model.ApplicationError
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&output))
	assert.Equal(t, AppErrorTypeStats, output.Type)
}

func TestGetStatistics_Ko_otherError(t *testing.T) {
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)

	stats := mocks.NewFizzBuzzStats(t)

	stats.On("Stats", req.Context()).Return(model.FizzBuzzStatisticsOutput{}, errors.New("dummy"))

	fbs := FizzBuzzServer{
		Stats: stats,
	}

	fbs.GetStatisticsHandler(resp, req)
	assert.Equal(t, http.StatusInternalServerError, resp.Result().StatusCode)
	var output model.ApplicationError
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&output))
	assert.Equal(t, AppErrorTypeStats, output.Type)
}
