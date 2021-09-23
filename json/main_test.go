package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonMarshalMap(t *testing.T) {
	expect := "{\n  \"apiVersion\": \"1.0\",\n  \"hase\": \"007\"\n}"
	bytes, err := MarshalMap("hase", "007")
	assert.Equal(t, expect, string(bytes))
	assert.Nil(t, err)
}
