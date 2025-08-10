package logger

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type testWriter struct {
	store []byte
}

func (tw *testWriter) Write(p []byte) (n int, err error) {
	tw.store = append(tw.store, p...)
	return len(p), nil
}
func (tw *testWriter) Sync() error {
	return nil
}

func TestInitialize(t *testing.T) {
	tests := []struct {
		name      string
		level     string
		expectErr bool
	}{
		{
			name:      "valid level",
			level:     "info",
			expectErr: false,
		},
		{
			name:      "invalid level fallback",
			level:     "invalid",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.RemoveAll("./logs")
			err := Initialize(tt.level)
			if (err != nil) != tt.expectErr {
				t.Fatalf("Initialize() error = %v, expectErr %v", err, tt.expectErr)
			}
			if _, err := os.Stat("./logs"); os.IsNotExist(err) {
				t.Error("logs directory was not created")
			}
			f, err := os.OpenFile("./logs/app.log", os.O_RDONLY|os.O_CREATE, 0644)
			if err != nil {
				t.Errorf("Failed to open log file: %v", err)
			}
			f.Close()
		})
	}
}

func TestRequestLogger(t *testing.T) {
	tw := &testWriter{}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(tw),
		zapcore.DebugLevel,
	)
	Log = zap.New(core)

	tests := []struct {
		name           string
		method         string
		url            string
		statusToReturn int
	}{
		{"GET 200", http.MethodGet, "/test", 200},
		{"POST 404", http.MethodPost, "/notfound", 404},
		{"PUT 500", http.MethodPut, "/error", 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tw.store = nil

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusToReturn)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.url, nil)
			req.RemoteAddr = "1.2.3.4:5678"

			RequestLogger(handler).ServeHTTP(w, req)

			if w.Result().StatusCode != tt.statusToReturn {
				t.Errorf("Expected status %d but got %d", tt.statusToReturn, w.Result().StatusCode)
			}

			logData := tw.store
			if len(logData) == 0 {
				t.Fatal("No log data captured")
			}

			logLines := strings.Split(strings.TrimSpace(string(logData)), "\n")
			var entry map[string]interface{}
			err := json.Unmarshal([]byte(logLines[0]), &entry)
			if err != nil {
				t.Fatalf("Failed to unmarshal log JSON: %v", err)
			}

			if entry["msg"] != "request logged" {
				t.Errorf("Expected msg 'request logged', got %v", entry["msg"])
			}
			if entry["method"] != tt.method {
				t.Errorf("Expected method %s, got %v", tt.method, entry["method"])
			}
			if entry["url"] != tt.url {
				t.Errorf("Expected url %s, got %v", tt.url, entry["url"])
			}
			if int(entry["status"].(float64)) != tt.statusToReturn {
				t.Errorf("Expected status %d, got %v", tt.statusToReturn, entry["status"])
			}
			if entry["remote_addr"] != req.RemoteAddr {
				t.Errorf("Expected remote_addr %s, got %v", req.RemoteAddr, entry["remote_addr"])
			}
			durVal, ok := entry["duration"]
			if !ok {
				t.Errorf("Expected duration field, got none")
			} else {
				durFloat, ok := durVal.(float64)
				if !ok {
					t.Errorf("Expected duration as float64, got %T", durVal)
				} else if durFloat <= 0 {
					t.Errorf("Expected positive duration, got %v", durFloat)
				}
			}
		})
	}
}
