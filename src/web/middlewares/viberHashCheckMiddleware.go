package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"io/ioutil"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// ViberHashCheck gets middleware for validating Viber callbacks
func ViberHashCheck(apiKey string) echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:X-Viber-Content-Signature",
		Validator: func(hash string, e echo.Context) (bool, error) {
			mac := hmac.New(sha256.New, []byte(apiKey))
			b, err := ioutil.ReadAll(e.Request().Body)
			if err != nil {
				return false, err
			}
			mac.Write(b)
			expectedHash := mac.Sum(nil)
			return hmac.Equal([]byte(hash), expectedHash), nil
		},
	})
}
