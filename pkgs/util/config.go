package util

import (
	"github.com/spf13/viper"
)

// Config defines the application config.
type Config struct {
	CloudFlareAPIToken    string `mapstructure:"CLOUDFLARE_API_TOKEN"`
	CloudFlareZoneID      string `mapstructure:"CLOUDFLARE_ZONE_ID"`
	CloudFlareEndpointURL string `mapstructure:"CLOUDFLARE_ENDPOINT_URL"`
	FrequencySeconds      int32  `mapstructure:"FREQUENCY_SECONDS"`
	PeriodSeconds         int32  `mapstructure:"PERIOD_SECONDS"`
	MetricsNamespace      string `mapstructure:"METRICS_NAMESPACE"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("local")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
