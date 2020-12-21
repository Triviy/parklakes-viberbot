package middlewares

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/triviy/parklakes-viberbot/infrastructure/logger"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// CustomLogger integrates with logrus
func CustomLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			logRequest(ctx)

			start := time.Now()
			if err = next(ctx); err != nil {
				ctx.Error(err)
			}
			stop := time.Now()

			logResponse(ctx, err, stop.Sub(start))

			return nil
		}
	}
}

func logRequest(ctx echo.Context) {
	req := ctx.Request()
	var body string
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Error(errors.Wrap(err, "reading from request body failed"))
	} else {
		body = string(b)
	}
	defer func() {
		req.Body.Close()
		req.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	}()
	logger.InfoDetailed(requestLogEntry{
		Method:        req.Method,
		URI:           req.RequestURI,
		Path:          req.URL.Path,
		RemoteIP:      ctx.RealIP(),
		Host:          req.Host,
		Protocol:      req.Proto,
		Referer:       req.Referer(),
		UserAgent:     req.UserAgent(),
		RequestID:     req.Header.Get(echo.HeaderXRequestID),
		Body:          body,
		ContentLength: req.ContentLength,
	}, "-- Start request")
}

func logResponse(ctx echo.Context, err error, latency time.Duration) {
	req := ctx.Request()
	res := ctx.Response()
	logger.InfoDetailed(responseLogEntry{
		Method:     req.Method,
		URI:        req.RequestURI,
		Path:       req.URL.Path,
		StatusCode: res.Status,
		Error:      getErrorText(err),
		RequestID:  res.Header().Get(echo.HeaderXRequestID),
		StackTrace: getStackTrace(err),
		Latency:    latency.String(),
	}, "-- End request")
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
