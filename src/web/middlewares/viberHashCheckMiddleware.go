package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"
)

// ViberHashCheck gets middleware for validating Viber callbacks
func ViberHashCheck(apiKey string) echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:X-Viber-Content-Signature",
		Validator: func(actualHash string, e echo.Context) (bool, error) {
			mac := hmac.New(sha256.New, []byte(apiKey))
			body, err := e.Request().GetBody()
			if err != nil {
				return false, errors.Wrap(err, "getting body from request failed")
			}
			b, err := ioutil.ReadAll(body)
			if err != nil {
				return false, errors.Wrap(err, "getting bytes from request body failed")
			}
			mac.Write(b)
			expectedHash := hex.EncodeToString(mac.Sum(nil))
			return hmac.Equal([]byte(actualHash), []byte(expectedHash)), nil
		},
	})
}
