package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	assert.Equal(t, "default", GetEnv("MY_ENV_VAR", "default"))
	os.Setenv("MY_ENV_VAR", "not_default")
	assert.Equal(t, "not_default", GetEnv("MY_ENV_VAR", "default"))
}

func TestIsTLSEnabled(t *testing.T) {
	assert.False(t, IsTLSEnabled("MY_TLS_ENV_VAR"))
	os.Setenv("MY_TLS_ENV_VAR", "true")
	assert.True(t, IsTLSEnabled("MY_TLS_ENV_VAR"))
	os.Setenv("MY_TLS_ENV_VAR", "TRUE")
	assert.True(t, IsTLSEnabled("MY_TLS_ENV_VAR"))
}
