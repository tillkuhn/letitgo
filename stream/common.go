package stream

import "github.com/rs/zerolog"

// LoggerWrapper wraps zerolog logger so we can used it as logger in kafka-go ReaderConfig
// Example:
// 		r := kafka.NewReader(kafka.ReaderConfig{
// 			Logger:      LoggerWrapper{delegate: k.logger},
// 		})
type LoggerWrapper struct {
	delegate zerolog.Logger
}

func (l LoggerWrapper) Printf(format string, v ...interface{}) {
	l.delegate.Trace().Msgf(format, v...)
}
