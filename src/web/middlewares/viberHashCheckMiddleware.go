package middlewares

import (
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
			log.Info("In KeyAuthWithConfig")
			mac := hmac.New(sha256.New, []byte(apiKey))
			log.Info("In KeyAuthWithConfig.Request.GetBody")
			body, err := e.Request().GetBody()
			defer body.Close()
			log.Info("After KeyAuthWithConfig.Request.GetBody")
			if err != nil {
				return false, errors.Wrap(err, "getting body from request failed")
			}
			log.Info("In KeyAuthWithConfig.ReadAll")
			b, err := ioutil.ReadAll(body)
			if err != nil {
				return false, errors.Wrap(err, "getting bytes from request body failed")
			}
			log.Info("In KeyAuthWithConfig.Write")
			mac.Write(b)
			log.Info("In KeyAuthWithConfig.EncodeToString")
			expectedHash := hex.EncodeToString(mac.Sum(nil))
			log.Info("In KeyAuthWithConfig.hmac.Equal")
			return hmac.Equal([]byte(actualHash), []byte(expectedHash)), nil
		},
	})
}
