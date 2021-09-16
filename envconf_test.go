package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	e := map[string]string {
		"T_AWS_REGION": "eu-mel-1",
		"T_SNS_ENABLED": "true",
		"T_MAX_ATTEMPTS": "42",
		"T_BROKERS": "localhost:44,localhost:55,localhost:66",
		"T_COLOR_CODES": "red:1,green:2,blue:3",
	}
	for k,v := range e {
		if err:=os.Setenv(k,v); err != nil {
			t.Fatal(err)
		}
	}
	c := ProcessEnv("t")
	assert.Equal(t, "eu-mel-1",c.AWSRegion)
	assert.Equal(t, true,c.SNSEnabled)
	assert.Equal(t, 42,c.MaxAttempts)
	assert.Equal(t, []string{"localhost:44","localhost:55","localhost:66"},c.Brokers)
	assert.Equal(t, 3,len(c.ColorCodes))

}
