package json

import (
	"github.com/tillkuhn/letitgo/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonMarshalMap(t *testing.T) {
	expect := "{\n  \"apiVersion\": \"1.0\",\n  \"hase\": \"007\"\n}"
	bytes, err := MarshalMapNicely("hase", "007")
	assert.Equal(t, expect, string(bytes))
	assert.Nil(t, err)
}

func TestMarshall(t *testing.T) {
	v := GenericUnmarshal[types.Bike](`{"Brand": "Merida"}`)
	assert.Equal(t, "Merida", v.Brand)
}
