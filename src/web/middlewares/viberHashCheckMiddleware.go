package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
)

// ViberHashCheck gets middleware for validating Viber callbacks
func ViberHashCheck(apiKey string) echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:X-Viber-Content-Signature",
		Validator: func(hash string, e echo.Context) (bool, error) {
			logrus.Infof("Got hash in request: %s", hash)

			mac := hmac.New(sha256.New, []byte(apiKey))
			b, err := ioutil.ReadAll(e.Request().Body)
			if err != nil {
				logrus.Infof("Error while reading body: %v", err)
				return false, err
			}
			mac.Write(b)
			expectedHash := mac.Sum(nil)
			logrus.Infof("Body: %v", string(b))
			logrus.Infof("Expected hash: %s", string(expectedHash))
			logrus.Infof("Expected encoded hash: %s", hex.EncodeToString(expectedHash))

			return hmac.Equal([]byte(hash), expectedHash), nil
		},
	})
}
