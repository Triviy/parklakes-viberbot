package middlewares

import (
	"fmt"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// CustomLogger integrates with logrus
func CustomLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			req := ctx.Request()
			res := ctx.Response()

			log.WithField("details", requestLogEntry{
				Method:        req.Method,
				URI:           req.RequestURI,
				Path:          req.URL.Path,
				RemoteIP:      ctx.RealIP(),
				Host:          req.Host,
				Protocol:      req.Proto,
				Referer:       req.Referer(),
				UserAgent:     req.UserAgent(),
				RequestID:     req.Header.Get(echo.HeaderXRequestID),
				ContentLength: req.ContentLength,
			}).Info("-- Start request")

			start := time.Now()
			if err = next(ctx); err != nil {
				ctx.Error(err)
			}
			stop := time.Now()
			// add some default fields to the logger ~ on all messages
			log.WithField("details", responseLogEntry{
				Method:     req.Method,
				URI:        req.RequestURI,
				Path:       req.URL.Path,
				StatusCode: res.Status,
				Error:      getErrorText(err),
				RequestID:  res.Header().Get(echo.HeaderXRequestID),
				StackTrace: getStackTrace(err),
				Latency:    stop.Sub(start).String(),
			}).Info("-- End request")
			return nil
		}
	}
}

type requestLogEntry struct {
	Method        string `json:"method"`
	URI           string `json:"uri"`
	Path          string `json:"path"`
	RemoteIP      string `json:"remoteIp"`
	Host          string `json:"host"`
	Protocol      string `json:"protocol"`
	Referer       string `json:"referer"`
	UserAgent     string `json:"userAgent"`
	RequestID     string `json:"requestId"`
	ContentLength int64  `json:"contentLength"`
}

type responseLogEntry struct {
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
	Method     string `json:"method"`
	URI        string `json:"uri"`
	Path       string `json:"path"`
	RequestID  string `json:"requestId"`
	Latency    string `json:"latency"`
	StackTrace string `json:"stackTrace"`
}

func getErrorText(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func getStackTrace(err error) string {
	if err == nil {
		return ""
	}

	if pkgErr, ok := err.(stackTracer); ok {
		var sb strings.Builder
		for _, f := range pkgErr.StackTrace() {
			sb.WriteString(fmt.Sprintf("%+s:%d\n", f, f))
		}
		return sb.String()
	}
	return err.Error()
}
