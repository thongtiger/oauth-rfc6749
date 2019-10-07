package main

import (
	"net/http"

	"github.com/thongtiger/oauth-rfc6749/auth"
	"github.com/thongtiger/oauth-rfc6749/handle"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.Match([]string{http.MethodGet, http.MethodPost}, "/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/oauth2/token", handle.TokenHandle)

	e.GET("/protected", func(c echo.Context) error {
		return c.String(http.StatusOK, "allow protected")
	}, auth.JWTMiddleware())

	e.Logger.Fatal(e.Start(":1323"))
}
