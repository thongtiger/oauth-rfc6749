package handle

import (
	"jwt-refresh-token/auth"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

const accessTokenDuration = time.Duration(time.Minute * 15)
const refreshTokenDuration = time.Duration(time.Minute * 30)

func TokenHandle(c echo.Context) (err error) {
	body := new(auth.Oauth2)
	if err = c.Bind(body); err != nil {
		return c.JSON(http.StatusUnsupportedMediaType, echo.Map{})
	}
	switch body.GrantType {
	case "password":
		/*
			- validateUser <- user
			- newRefreshToken <- refresh_token
			- newAccessToken <- access_token
		*/
		if ok, user := auth.ValidateUser(body.Username, body.Password); ok {
			GenerateTK(c, user)
		}

	case "refresh_token":
		/*
			- check exists refreshToken
			- get refreshToken info
			- replace old refreshToken
			- response new access_token,refresh_token
		*/
		if ok, claim := auth.ValidateRefreshToken(body.RefreshToken); ok {
			log.Println(claim)
			// GenerateTK(c, user)
		}
	}
	return c.JSON(http.StatusUnauthorized, echo.Map{})
}

//GenerateTK generate response token
func GenerateTK(c echo.Context, user auth.User) (err error) {
	accessToken, err := auth.NewToken(user.ID.Hex(), accessTokenDuration, "access_token", user.Role, user.Scope...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "can't generate access_token", "err": err.Error()})
	}
	refreshToken, err := auth.NewToken(user.ID.Hex(), refreshTokenDuration, "refresh_token", user.Role, user.Scope...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "can't generate refresh_token"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"client_id":     user.ID,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "bearer",
		"expires_in":    int64(accessTokenDuration.Seconds()),
		"scope":         user.Scope,
	})
}
