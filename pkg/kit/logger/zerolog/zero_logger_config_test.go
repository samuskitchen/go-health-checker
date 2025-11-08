package zerolog

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

// newTestLogger creates a zerolog.Logger that writes to the passed buffer.
// We use a colorless ConsoleWriter to make asserting text easier.
func newTestLogger(output *bytes.Buffer) zerolog.Logger {
	writer := zerolog.ConsoleWriter{
		Out:        output,
		NoColor:    true,
		PartsOrder: []string{"time", "level", "message"},
	}
	return zerolog.New(writer).With().Timestamp().Logger()
}

func TestLogWithConfig(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name       string
		skipper    func(echo.Context) bool
		expectNext bool // whether the next handler should be executed
		expectLog  bool // whether there should be log output
		fieldMap   map[string]string
	}{
		{
			name:       "skipper=true ⇒ llama next y sin log",
			skipper:    func(echo.Context) bool { return true },
			expectNext: true,  // next is called inside the skipper block
			expectLog:  false, // does not log anything
			fieldMap:   map[string]string{"dummy": "@dummy"},
		},
		{
			name:       "skipper=false ⇒ llama next y sí loggea",
			skipper:    mw.DefaultSkipper,
			expectNext: true, // MapFields also invokes next
			expectLog:  true, // must register "handle request"
			fieldMap:   map[string]string{"uri": "@uri"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/foo?bar=baz", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			buf := &bytes.Buffer{}
			cfg := Config{
				Skipper:  tt.skipper,
				Logger:   newTestLogger(buf),
				FieldMap: tt.fieldMap,
			}

			called := false
			next := func(c echo.Context) error {
				called = true
				return c.String(http.StatusTeapot, "next-called")
			}

			handler := LogWithConfig(cfg)(next)
			if err := handler(ctx); err != nil {
				t.Fatalf("middleware returned error: %v", err)
			}

			// 1) We check if next was called
			if called != tt.expectNext {
				t.Errorf("expected next called = %v, got %v", tt.expectNext, called)
			}

			out := buf.String()

			// 2) We check if there was a log
			if tt.expectLog {
				if !strings.Contains(out, "handle request") {
					t.Error("expected log output to contain message 'handle request'")
				}
				for field := range tt.fieldMap {
					if !strings.Contains(out, field) {
						t.Errorf("expected log output to contain field key '%s'", field)
					}
				}
			} else {
				if out != "" {
					t.Errorf("expected no log output, but got %q", out)
				}
			}
		})
	}
}
