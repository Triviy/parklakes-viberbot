package middlewares

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// APIKeyAuth gets middleware for validating API by APIKey
func APIKeyAuth(apiKey string) echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:X-API-KEY",
		Validator: func(s string, e echo.Context) (bool, error) {
			return apiKey == s, nil
		},
	})
}
