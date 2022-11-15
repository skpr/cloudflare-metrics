package util

import (
	"github.com/spf13/viper"
)

// Config defines the application config.
type Config struct {
	CloudFlareAPIToken    string   `mapstructure:"CLOUDFLARE_API_TOKEN"`
	CloudFlareZoneTags    []string `mapstructure:"CLOUDFLARE_ZONE_TAGS"`
	CloudFlareEndpointURL string   `mapstructure:"CLOUDFLARE_ENDPOINT_URL"`
	PeriodSeconds         int32    `mapstructure:"PERIOD_SECONDS"`
	MetricsNamespace      string   `mapstructure:"METRICS_NAMESPACE"`
}

// Validate validates the config.
func (c Config) Validate() []string {
	var errors []string
	if c.CloudFlareAPIToken == "" {
		errors = append(errors, "CLOUDFLARE_API_TOKEN is a required variable")
	}
	if len(c.CloudFlareZoneTags) == 0 {
		errors = append(errors, "CLOUDFLARE_ZONE_TAGS is a required variable")
	}
	if c.CloudFlareEndpointURL == "" {
		errors = append(errors, "CLOUDFLARE_ENDPOINT_URL is a required variable")
	}
	if c.PeriodSeconds < 60 {
		errors = append(errors, "PERIOD_SECONDS must be a valid integer greater than 60")
	}
	if c.MetricsNamespace == "" {
		errors = append(errors, "METRICS_NAMESPACE is a required variable")
	}
	return errors
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("local")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	_ = viper.ReadInConfig()
	err = viper.Unmarshal(&config)
	return
}
