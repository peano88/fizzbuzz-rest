package server

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigureServer_NoTLS(t *testing.T) {
	tbs := FizzBuzzServer{}

	s, err := tbs.Configure()
	assert.NoError(t, err)
	assert.NotNil(t, s)
}

func TestConfigureServer_TLS(t *testing.T) {
	tbs := FizzBuzzServer{}

	os.Setenv(TLSEnvVar, "true")

	s, err := tbs.Configure()
	assert.NoError(t, err)
	assert.NotNil(t, s)

	// creating issues
	os.Setenv(insecureEnvVar, "boh")
	s, err = tbs.Configure()
	assert.Nil(t, s)
	assert.Error(t, err)
	os.Unsetenv(insecureEnvVar)

	os.Setenv(clientAuthTypeEnvVar, "boh")
	s, err = tbs.Configure()
	assert.Nil(t, s)
	assert.Error(t, err)
	os.Unsetenv(clientAuthTypeEnvVar)
}
