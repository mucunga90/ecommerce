package service

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type CustomClaims struct {
	Scope       string           `json:"scope,omitempty"`
	Scopes      []string         `json:"scp,omitempty"`
	Roles       []string         `json:"roles,omitempty"`
	Permissions []string         `json:"permissions,omitempty"`
	Audience    jwt.ClaimStrings `json:"aud,omitempty"`
	jwt.RegisteredClaims
}

type Verifier struct {
	jwks   *keyfunc.JWKS
	issuer string
	aud    string
}

func NewVerifier(ctx context.Context, issuer, jwksURL, audience string) (*Verifier, error) {
	if jwksURL == "" {
		if strings.Contains(issuer, "/realms/") {
			jwksURL = strings.TrimSuffix(issuer, "/") + "/protocol/openid-connect/certs"
		} else {
			jwksURL = strings.TrimSuffix(issuer, "/") + "/.well-known/jwks.json"
		}
	}

	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshInterval:     time.Minute * 5,
		RefreshErrorHandler: func(err error) {}, // you can log
		RefreshUnknownKID:   true,
		RefreshTimeout:      time.Second * 10,
		Ctx:                 ctx,
	})
	if err != nil {
		return nil, err
	}

	return &Verifier{
		jwks:   jwks,
		issuer: issuer,
		aud:    audience,
	}, nil
}

func (v *Verifier) EchoJWTMiddleware() echo.MiddlewareFunc {
	cfg := echojwt.Config{
		ContextKey:  "user",
		TokenLookup: "header:Authorization:Bearer ",
		KeyFunc: func(token *jwt.Token) (interface{}, error) {
			return v.jwks.Keyfunc(token)
		},
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(CustomClaims)
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or missing token"})
		},
	}
	return echojwt.WithConfig(cfg)
}

func (v *Verifier) RequireValidClaims(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tok, ok := c.Get("user").(*jwt.Token)
		if !ok || tok == nil || !tok.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
		}
		claims, ok := tok.Claims.(*CustomClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid claims"})
		}

		if claims.Issuer != v.issuer {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "bad issuer"})
		}
		if v.aud != "" && !containsAudience(claims.Audience, v.aud) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "bad audience"})
		}

		scopes := map[string]struct{}{}
		for _, s := range strings.Fields(claims.Scope) {
			scopes[s] = struct{}{}
		}
		for _, s := range claims.Scopes {
			scopes[s] = struct{}{}
		}
		for _, s := range claims.Permissions {
			scopes[s] = struct{}{}
		}

		c.Set("scopes", scopes)
		c.Set("roles", claims.Roles)
		return next(c)
	}
}

func containsAudience(audiences jwt.ClaimStrings, target string) bool {
	for _, aud := range audiences {
		if aud == target {
			return true
		}
	}
	return false
}
