package zerolog

import (
	"github.com/samuskitchen/go-health-checker/pkg/kit/logger"

	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config defines the config for ZeroLog middleware.
type Config struct {
	// FieldMap set a list of fields with tags
	//
	//  to construct the logger fields.
	//
	// - @id (Request ID)
	// - @remote_ip
	// - @uri
	// - @host
	// - @method
	// - @path
	// - @protocol
	// - @referer
	// - @user_agent
	// - @status
	// - @error
	// - @latency (In nanoseconds)
	// - @latency_human (Human readable)
	// - @bytes_in (Bytes received)
	// - @bytes_out (Bytes sent)
	// - @header:<NAME>
	// - @query:<NAME>
	// - @form:<NAME>
	// - @cookie:<NAME>
	FieldMap map[string]string

	// Logger it is a zero-log logger
	Logger zerolog.Logger

	// Skipper defines a function to skip middleware.
	Skipper mw.Skipper
}

// DefaultLogConfig is the default ZeroLog middleware config.
var DefaultLogConfig = Config{
	FieldMap: logger.DefaultFields,
	Logger:   log.Logger,
	Skipper:  mw.DefaultSkipper,
}

// LogWithConfig returns ZeroLog middleware with config.
// See: `ZeroLog()`.
func LogWithConfig(cfg Config) echo.MiddlewareFunc {
	// Defaults
	if cfg.Skipper == nil {
		cfg.Skipper = DefaultLogConfig.Skipper
	}

	if len(cfg.FieldMap) == 0 {
		cfg.FieldMap = DefaultLogConfig.FieldMap
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			if cfg.Skipper(ctx) {
				return next(ctx)
			}

			logFields, err := logger.MapFields(ctx, next, cfg.FieldMap)

			cfg.Logger.Info().
				Fields(logFields).
				Msg("handle request")

			return
		}
	}
}
