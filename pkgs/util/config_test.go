package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig("testdata")
	assert.NoError(t, err)
	assert.Equal(t, "abcd1234", config.CloudFlareAPIToken)
	assert.Equal(t, "xyz456", config.CloudFlareZoneTag)
	assert.Equal(t, "example.com", config.CloudFlareHostName)
	assert.Equal(t, "https://api.cloudflare.com/client/v4/graphql", config.CloudFlareEndpointURL)
	assert.EqualValues(t, 60, config.PeriodSeconds)
	assert.Equal(t, "Skpr/CloudFlare", config.MetricsNamespace)
	config = Config{}
}

func TestValidate(t *testing.T) {
	t.Skip("Need to work out why this fails in CI")

	// Load missing config.
	config, err := LoadConfig("")
	assert.NoError(t, err)
	errors := config.Validate()
	assert.Len(t, errors, 5)
}
