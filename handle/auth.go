package handle

import (
	"jwt-refresh-token/auth"
	"jwt-refresh-token/redis"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

const (
	accessTokenDuration  = time.Duration(time.Minute * 15)
	refreshTokenDuration = time.Duration(time.Minute * 30)
)

// TokenHandle route
func TokenHandle(c echo.Context) (err error) {
	body := new(auth.Oauth2)
	if err = c.Bind(body); err != nil {
		return c.JSON(http.StatusUnsupportedMediaType, echo.Map{})
	}
	switch body.GrantType {
	case "password":
		if ok, user := auth.ValidateUser(body.Username, body.Password); ok {
			// generate token
			return GenerateTK(c, user)
		}
	case "refresh_token":
		// validate
		tokenValid, claim := auth.ValidateRefreshToken(body.RefreshToken)
		if !tokenValid {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "bad refresh_token"})
		}
		// exists
		if exist := redis.Exists(claim.ID, body.RefreshToken); !exist {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "refresh_token does not exist or expired"})
		}
		// generate token
		user := auth.User{ID: bson.ObjectIdHex(claim.ID), Role: claim.Role, Scope: claim.Scope}
		return GenerateTK(c, user)
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