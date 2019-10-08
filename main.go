package main

import (
	"github.com/labstack/echo/middleware"
	"net/http"

	"github.com/thongtiger/oauth-rfc6749/auth"
	"github.com/thongtiger/oauth-rfc6749/handle"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentLength, echo.HeaderAcceptEncoding, echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderXRequestedWith},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
		MaxAge:       3600,
	}))

	e.Match([]string{http.MethodGet, http.MethodPost}, "/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/oauth2/token", handle.TokenHandle)

	e.GET("/protected", func(c echo.Context) error {
		return c.String(http.StatusOK, "allow protected")
	}, auth.JWTMiddleware())
	e.GET("/401", func(c echo.Context) error {
		return echo.ErrUnauthorized
	})
	e.Logger.Fatal(e.Start(":1323"))
}
