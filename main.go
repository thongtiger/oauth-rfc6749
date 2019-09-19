package main

import (
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

	e.Logger.Fatal(e.Start(":1323"))
}
