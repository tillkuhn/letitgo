package main

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

// Config Specification
type Config struct {
	AWSRegion   string   `default:"eu-central-1" required:"false" desc:"AWS Region" split_words:"true"`
	SNSEnabled  bool     `default:"false" required:"false" desc:"Enable SNS" split_words:"true"`
	MaxAttempts int      `default:"4" required:"true" desc:"Max attempts" split_words:"true"`
	Brokers     []string `required:"false" desc:"List of brokers" split_words:"true"`
	ColorCodes  map[string]int `split_words:"true"`
}

// ProcessEnv configures the app with https://github.com/kelseyhightower/envconfig
func ProcessEnv(appPrefix string) *Config {
	var config Config
	// populates the appPrefix struct based on environment variables
	err := envconfig.Process(appPrefix, &config)
	if err != nil {
		log.Fatalf("Cannot process envconfig: %v", err.Error())
	}
	return &config
}

