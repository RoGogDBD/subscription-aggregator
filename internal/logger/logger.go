package logger

import (
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func Initialize(level string) error {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{
		"./logs/app.log",
		"stdout",
	}

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if err := os.MkdirAll("./logs", 0755); err != nil {
		return err
	}

	var err error
	Log, err = config.Build()
	if err != nil {
		return err
	}
	return nil
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sr := &statusRecorder{ResponseWriter: w, status: 200}
		next.ServeHTTP(sr, r)
		duration := time.Since(start)

		Log.Info("request logged",
			zap.String("method", r.Method),
			zap.String("url", r.RequestURI),
			zap.Int("status", sr.status),
			zap.Duration("duration", duration),
			zap.String("remote_addr", r.RemoteAddr),
		)
	})
}
