package ctx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, requestIdKey, "12345")
	ctx = context.WithValue(ctx, userIdKey, "it's me")
	err := DoWithContext(ctx)
	assert.Nil(t, err)

	ctx = context.Background()
	err = DoWithContext(ctx)
	assert.NotNil(t, err)
}
