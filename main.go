package main

import (
	"jwt-refresh-token/auth"
	"jwt-refresh-token/handle"
	"net/http"

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
