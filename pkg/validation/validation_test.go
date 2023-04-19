package validation

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/peano88/fizzbuzz-rest/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFizzBuzzValidator_OK(t *testing.T) {

	v := NewFizzBuzzValidator()

	r := httptest.NewRequest(http.MethodGet, "http://example.com?int1=2&int2=3&limit=7&start=4&str1=fzz&str2=bzz", nil)

	ctx, err := v.RunValidations(r)
	require.NoError(t, err)
	genericInput := ctx.Value(model.InputKey)
	require.NotNil(t, genericInput)
	realInput, ok := genericInput.(model.FizzBuzzInput)
	assert.True(t, ok)
	assert.Equal(t, 2, realInput.Int1)
	assert.Equal(t, 3, realInput.Int2)
	assert.Equal(t, 7, realInput.Limit)
	assert.Equal(t, 4, realInput.Start)
	assert.Equal(t, "fzz", realInput.Str1)
	assert.Equal(t, "bzz", realInput.Str2)
}

func TestFizzBuzzValidator_KO(t *testing.T) {
	v := NewFizzBuzzValidator()

	r := httptest.NewRequest(http.MethodGet, "http://example.com?&int2=3&limit=7&str1=fzz&str2=bzz", nil)
	ctx, err := v.RunValidations(r)
	assert.Error(t, err)
	assert.Nil(t, ctx)

	r = httptest.NewRequest(http.MethodGet, "http://example.com?int1=-1&int2=3&limit=7&str1=fzz&str2=bzz", nil)
	ctx, err = v.RunValidations(r)
	assert.Error(t, err)
	assert.Nil(t, ctx)

	r = httptest.NewRequest(http.MethodGet, "http://example.com?int2=22222222222222222222222222222222222222222222222222222&int1=3&limit=7&str1=fzz&str2=bzz", nil)
	ctx, err = v.RunValidations(r)
	assert.Error(t, err)
	assert.Nil(t, ctx)

	r = httptest.NewRequest(http.MethodGet, "http://example.com?int1=2&int2=3&limit=seven&str1=fzz&str2=bzz", nil)
	ctx, err = v.RunValidations(r)
	assert.Error(t, err)
	assert.Nil(t, ctx)

	r = httptest.NewRequest(http.MethodGet, "http://example.com?int1=2&int2=3&limit=7&str1=fzz", nil)
	ctx, err = v.RunValidations(r)
	assert.Error(t, err)
	assert.Nil(t, ctx)

	r = httptest.NewRequest(http.MethodGet, "http://example.com?int1=2&int2=3&limit=7&str1=f-zz&str2=bzz", nil)
	ctx, err = v.RunValidations(r)
	assert.Error(t, err)
	assert.Nil(t, ctx)

	r = httptest.NewRequest(http.MethodGet, "http://example.com?int1=2&int2=3&limit=7&start=blah&str1=fzz&str2=bzz", nil)
	ctx, err = v.RunValidations(r)
	assert.Error(t, err)
	assert.Nil(t, ctx)

	r = httptest.NewRequest(http.MethodGet, "http://example.com?int1=2&int2=3&str1=fzz&str2=bzz", nil)
	ctx, err = v.RunValidations(r)
	assert.Error(t, err)
	assert.Nil(t, ctx)
}

func TestValidationError(t *testing.T) {
	errA := errors.New("error A")
	valErr := ValidationError{
		err:        errA,
		parameter:  "donald",
		constraint: "duck",
	}

	assert.NotEmpty(t, valErr.Error())
	assert.True(t, errors.Is(valErr, errA))
	assert.Equal(t, "duck", valErr.Constraint())
}
