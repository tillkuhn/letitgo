package ctx

import (
	"context"
	"errors"
	"log"
)

// https://blog.gopheracademy.com/advent-2016/context-logging/
type contextKey int

const (
	requestIdKey contextKey = iota
	userIdKey

)


func DoWithContext(ctx context.Context) error {
	requestId := ctx.Value(requestIdKey)
	if requestId == nil {
		return errors.New("requestId not set")
	}
	log.Printf("requestId %v from %v",requestId,ctx.Value(userIdKey))
	return nil
}
