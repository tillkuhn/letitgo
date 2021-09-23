package json

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJsonMarshalMap(t *testing.T) {
	expect := "{\n  \"apiVersion\": \"1.0\",\n  \"hase\": \"007\"\n}"
	bytes,err := MarshalMap("hase","007")
	assert.Equal(t, expect,string(bytes))
	assert.Nil(t, err)
}
