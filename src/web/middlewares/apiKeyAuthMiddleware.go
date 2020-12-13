package middlewares

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// GetAPIKeyAuthMiddleware gets middleware for validating API by APIKey
func GetAPIKeyAuthMiddleware(apiKey string) echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:X-API-KEY",
		Validator: func(s string, e echo.Context) (bool, error) {
			return apiKey == s, nil
		},
	})
}

// GetViberAPIKeyAuthMiddleware gets middleware for validating Viber callbacks
func GetViberAPIKeyAuthMiddleware(apiKey string) echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:X-Viber-Auth-Token",
		Validator: func(s string, e echo.Context) (bool, error) {
			return apiKey == s, nil
		},
	})
}
