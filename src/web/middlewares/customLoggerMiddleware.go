package middlewares

import (
	"fmt"
	"net/http"
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
			start := time.Now()
			if err = next(ctx); err != nil {
				ctx.Error(err)
			}
			stop := time.Now()
			// add some default fields to the logger ~ on all messages
			log.WithField("request", logEntry{
				Method:        req.Method,
				URI:           req.RequestURI,
				Path:          req.URL.Path,
				RemoteIP:      ctx.RealIP(),
				Host:          req.Host,
				Protocol:      req.Proto,
				Referer:       req.Referer(),
				UserAgent:     req.UserAgent(),
				RequestID:     getRequestID(req, res),
				StatusCode:    res.Status,
				Error:         getErrorText(err),
				ContentLength: req.ContentLength,
				Latency:       stop.Sub(start).String(),
			}).Info("New request")
			return nil
		}
	}
}

type logEntry struct {
	Method        string `json:"method"`
	URI           string `json:"uri"`
	Path          string `json:"path"`
	RemoteIP      string `json:"remote_ip"`
	Host          string `json:"host"`
	Protocol      string `json:"protocol"`
	Referer       string `json:"referer"`
	UserAgent     string `json:"user_agent"`
	RequestID     string `json:"request_id"`
	StatusCode    int    `json:"status_code"`
	Error         string `json:"error"`
	ContentLength int64  `json:"content_length"`
	Latency       string `json:"latency"`
}

func getRequestID(req *http.Request, res *echo.Response) string {
	id := req.Header.Get(echo.HeaderXRequestID)
	if id == "" {
		id = res.Header().Get(echo.HeaderXRequestID)
	}
	return id
}

func getErrorText(err error) string {
	if err == nil {
		return ""
	}

	if pkgErr, ok := err.(stackTracer); ok {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%s\n", err.Error()))
		for _, f := range pkgErr.StackTrace() {
			sb.WriteString(fmt.Sprintf("%+s:%d\n", f, f))
		}
		return sb.String()
	}
	return err.Error()
}
