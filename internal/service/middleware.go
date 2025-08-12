package service

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RequireScopes(required ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			raw := c.Get("scopes")
			scopes, ok := raw.(map[string]struct{})
			if !ok {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "missing scopes"})
			}
			for _, r := range required {
				if _, ok := scopes[r]; !ok {
					return c.JSON(http.StatusForbidden, map[string]string{"error": "insufficient_scope"})
				}
			}
			return next(c)
		}
	}
}

func RequireAnyRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			raw := c.Get("roles")
			have, _ := raw.([]string)
			for _, need := range roles {
				for _, r := range have {
					if r == need {
						return next(c)
					}
				}
			}
			return c.JSON(http.StatusForbidden, map[string]string{"error": "insufficient_role"})
		}
	}
}
