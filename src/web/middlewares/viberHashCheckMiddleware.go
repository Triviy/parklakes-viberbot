package middlewares

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// ViberHashCheck gets middleware for validating Viber callbacks
func ViberHashCheck(apiKey string) echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:X-Viber-Content-Signature",
		Validator: func(actualHash string, e echo.Context) (bool, error) {
			mac := hmac.New(sha256.New, []byte(apiKey))

			var buf bytes.Buffer
			tee := io.TeeReader(e.Request().Body, &buf)
			b, err := ioutil.ReadAll(tee)

			if err != nil {
				return false, err
			}
			mac.Write(b)
			expectedHash := hex.EncodeToString(mac.Sum(nil))
			return hmac.Equal([]byte(actualHash), []byte(expectedHash)), nil
		},
	})
}
