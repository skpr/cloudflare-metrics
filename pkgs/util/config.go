package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// Config defines the application config.
type Config struct {
	CloudFlareAPIToken    string            `mapstructure:"CLOUDFLARE_METRICS_API_TOKEN" validate:"required"`
	CloudFlareZoneTag     string            `mapstructure:"CLOUDFLARE_METRICS_ZONE_TAG" validate:"required"`
	CloudFlareHostNames   []string          `mapstructure:"CLOUDFLARE_METRICS_HOSTNAMES"`
	CloudFlareEndpointURL string            `mapstructure:"CLOUDFLARE_METRICS_ENDPOINT_URL"`
	Period                time.Duration     `mapstructure:"CLOUDFLARE_METRICS_PERIOD"`
	MetricsNamespace      string            `mapstructure:"CLOUDFLARE_METRICS_NAMESPACE" validate:"required"`
	ExtraDimensionsRaw    []string          `mapstructure:"CLOUDFLARE_METRICS_EXTRA_DIMENSIONS"`
	ExtraDimensions       map[string]string `mapstructure:"-"`
}

// GetDimensionsMap gets the dimension map.
func (c *Config) buildDimensionsMap() error {
	d := make(map[string]string)
	for _, c := range c.ExtraDimensionsRaw {
		parts := strings.Split(c, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid dimension format %s", c)
		}
		d[parts[0]] = parts[1]
	}
	c.ExtraDimensions = d
	return nil
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("local")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	viper.SetDefault("CLOUDFLARE_METRICS_PERIOD", time.Minute)
	viper.SetDefault("CLOUDFLARE_METRICS_ENDPOINT_URL", "https://api.cloudflare.com/client/v4/graphql")

	// We don't use a dotenv file in prod.
	_ = viper.ReadInConfig()

	var config Config
	err := viper.Unmarshal(&config)

	validate := validator.New()
	err = validate.Struct(&config)
	if err != nil {
		return Config{}, err
	}
	err = config.buildDimensionsMap()
	return config, err
}
