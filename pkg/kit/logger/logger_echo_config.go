package logger

import (
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// Log middlewares constants.
const (
	logID           = "@id"
	logRemoteIP     = "@remote_ip"
	logURI          = "@uri"
	logHost         = "@host"
	logMethod       = "@method"
	logPath         = "@path"
	logProtocol     = "@protocol"
	logReferer      = "@referer"
	logUserAgent    = "@user_agent"
	logStatus       = "@status"
	logError        = "@error"
	logLatency      = "@latency"
	logLatencyHuman = "@latency_human"
	logBytesIn      = "@bytes_in"
	logBytesOut     = "@bytes_out"
	logHeaderPrefix = "@header:"
	logQueryPrefix  = "@query:"
	logFormPrefix   = "@form:"
	logCookiePrefix = "@cookie:"
)

// DefaultFields for logger
var DefaultFields = map[string]string{
	"remote_ip": logRemoteIP,
	"uri":       logURI,
	"host":      logHost,
	"method":    logMethod,
	"status":    logStatus,
	"latency":   logLatency,
	"error":     logError,
}

// string to int base conversion.
const base = 10

// MapFields maps fields based on tag name.
func MapFields(echoContext echo.Context, handlerFunc echo.HandlerFunc, fieldsMap map[string]string) (
	map[string]interface{}, error,
) {
	logFields := map[string]interface{}{}
	start := time.Now()

	err := handlerFunc(echoContext)
	if err != nil {
		echoContext.Error(err)
	}

	elapsed := time.Since(start)
	tags := MapTags(echoContext, elapsed)

	if err != nil {
		tags[logError] = err
	}

	for index, tag := range fieldsMap {
		if tag == "" {
			continue
		}

		if value, ok := tags[tag]; ok {
			logFields[index] = value
			continue
		}

		switch {
		case strings.HasPrefix(tag, logHeaderPrefix):
			key := tag[len(logHeaderPrefix):]
			logFields[index] = echoContext.Request().Header.Get(key)

		case strings.HasPrefix(tag, logQueryPrefix):
			key := tag[len(logQueryPrefix):]
			logFields[index] = echoContext.QueryParam(key)

		case strings.HasPrefix(tag, logFormPrefix):
			key := tag[len(logFormPrefix):]
			logFields[index] = echoContext.FormValue(key)

		case strings.HasPrefix(tag, logCookiePrefix):
			key := tag[len(logCookiePrefix):]
			cookie, er := echoContext.Cookie(key)
			if er == nil {
				logFields[index] = cookie.Value
			}
		}
	}

	return logFields, err
}

// MapTags maps the log tags with its related data. Populate previously the
// key/value avoids the cyclomatic complexity of the log middlewares to
// identify each tag and value.
func MapTags(echoContext echo.Context, latency time.Duration) map[string]interface{} {
	tags := map[string]interface{}{}

	req := echoContext.Request()
	res := echoContext.Response()

	id := req.Header.Get(echo.HeaderXRequestID)
	if id == "" {
		id = res.Header().Get(echo.HeaderXRequestID)
	}

	tags[logID] = id
	tags[logRemoteIP] = echoContext.RealIP()
	tags[logURI] = req.RequestURI
	tags[logHost] = req.Host
	tags[logMethod] = req.Method

	path := req.URL.Path
	if path == "" {
		path = "/"
	}

	tags[logPath] = path
	tags[logProtocol] = req.Proto
	tags[logReferer] = req.Referer()
	tags[logUserAgent] = req.UserAgent()
	tags[logStatus] = res.Status
	tags[logLatency] = strconv.FormatInt(int64(latency), base)
	tags[logLatencyHuman] = latency.String()

	contentLength := req.Header.Get(echo.HeaderContentLength)
	if contentLength == "" {
		contentLength = "0"
	}

	tags[logBytesIn] = contentLength
	tags[logBytesOut] = strconv.FormatInt(res.Size, base)

	return tags
}
