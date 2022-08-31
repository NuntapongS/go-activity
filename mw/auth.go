package mw

import (
	"encoding/base64"
	"myapp/login"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func OnlyAdmin(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*login.JwtCustomClaims)
		if claims.Role != "ADMIN" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}
		return h(c)
	}
}

func SkipperAPIKey(c echo.Context) bool {
	authorization := c.Request().Header.Get("Authorization")
	if strings.TrimPrefix(authorization, "Bearer ") != "" {
		return false
	}

	apiKeyHeader := c.Request().Header.Get("API-KEY")
	apiKeyDec, _ := base64.StdEncoding.DecodeString(viper.GetString("api.key.public"))

	return apiKeyHeader == string(apiKeyDec)
}
