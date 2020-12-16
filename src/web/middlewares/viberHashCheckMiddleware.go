package middlewares

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
)

// ViberHashCheck gets middleware for validating Viber callbacks
func ViberHashCheck(apiKey string) echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:X-Viber-Content-Signature",
		Validator: func(actualHash string, e echo.Context) (bool, error) {
			req := e.Request()
			mac := hmac.New(sha256.New, []byte(apiKey))
			b, err := ioutil.ReadAll(req.Body)
			defer func() {
				req.Body.Close()
				log.Info(string(b))
				req.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()

			if err != nil {
				return false, errors.Wrap(err, "getting bytes from request body failed")
			}
			mac.Write(b)
			expectedHash := hex.EncodeToString(mac.Sum(nil))
			return hmac.Equal([]byte(actualHash), []byte(expectedHash)), nil
		},
	})
}
