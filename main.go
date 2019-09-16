package main

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

const accessTokenDuration = time.Duration(time.Minute * 15)
const refreshTokenDuration = time.Duration(time.Minute * 30)

type Oauth2 struct {
	Username     string `json:"username" form:"username" query:"username"`
	Password     string `json:"password" form:"password" query:"password"`
	GrantType    string `json:"grant_type" form:"grant_type" query:"grant_type"`
	RefreshToken string `json:"refresh_token" form:"refresh_token" query:"refresh_token"`
}

type User struct {
	ID             bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Role           string        `json:"role" bson:"role,omitempty"`
	Scope          []string      `json:"scope"`
	Username       string        `json:"username" bson:"username,omitempty"`
	Password       string        `json:"password" bson:"password,omitempty"`
	Name           string        `json:"name" bson:"name,omitempty"`
	CreateTime     time.Time     `json:"createTime" bson:"createTime"`
	LatestLoggedin time.Time     `json:"latestLoggedin" bson:"latestLoggedin"`
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/oauth2/token", func(c echo.Context) (err error) {
		body := new(Oauth2)
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
			if ok, user := ValidateUser(body.Username, body.Password); ok {
				GenerateTK(c, user)
			}

		case "refresh_token":
			/*
				- check exists refreshToken
				- get refreshToken info
				- replace old refreshToken
				- response new access_token,refresh_token
			*/
			if ok, claim := ValidateRefreshToken(body.RefreshToken); ok {
				log.Println(claim)
				// GenerateTK(c, user)
			}
		}

		return c.JSON(http.StatusUnauthorized, echo.Map{})
	})

	e.Logger.Fatal(e.Start(":3000"))
}

//GenerateTK generate response token
func GenerateTK(c echo.Context, user User) (err error) {
	accessToken, err := NewToken(user.ID.Hex(), accessTokenDuration, "access_token", user.Role, user.Scope...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "can't generate access_token", "err": err.Error()})
	}
	refreshToken, err := NewToken(user.ID.Hex(), refreshTokenDuration, "refresh_token", user.Role, user.Scope...)
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
