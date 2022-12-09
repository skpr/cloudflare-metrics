package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig("testdata")
	assert.NoError(t, err)
	assert.Equal(t, "abcd1234", config.CloudFlareAPIToken)
	assert.Equal(t, "xyz456", config.CloudFlareZoneTag)
	assert.Contains(t, config.CloudFlareHostNames, "example.com")
	assert.Equal(t, "https://api.cloudflare.com/client/v4/graphql", config.CloudFlareEndpointURL)
	assert.EqualValues(t, time.Minute, config.Period)
	assert.Equal(t, "Skpr/CloudFlare", config.MetricsNamespace)

	assert.Equal(t, "foo", config.ExtraDimensions["project"])
	assert.Equal(t, "bar", config.ExtraDimensions["env"])
}
