package auth

import (
	"log"
	"time"

	"github.com/thongtiger/oauth-rfc6749/redis"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/mgo.v2/bson"
)

const secretJWT = "secret"

// ValidateUser validates credentials of a potential user
func ValidateUser(username, password string) (bool, User) {
	if username == "joe" && password == "password" {
		return true, User{
			ID:       bson.NewObjectId(),
			Name:     "ioe",
			Username: "ioe",
			Role:     "emp",
			Scope:    []string{"1", "2"},
		}
	}
	return false, User{}
}
func ValidateRefreshToken(tokenString string) (bool, *TokenClaim) {

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretJWT), nil
	})
	if err != nil {
		return false, nil
	}
	if claims, ok := token.Claims.(*TokenClaim); ok && token.Valid {
		log.Printf("%v %v", claims.ID, claims.StandardClaims.ExpiresAt)
		return true, claims

	}
	return false, nil

}

func JWTMiddleware() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &TokenClaim{},
		SigningKey: []byte(secretJWT),
		ErrorHandler: func(err error) error {
			return echo.ErrUnauthorized
		},
		Skipper: func(c echo.Context) bool {
			// Skip authentication for and signup login requests
			if c.Path() == "/login" || c.Path() == "/signup" {
				return true
			}
			return false
		},
	}) //echo.HandlerFunc
}

func NewToken(id string, expiresIn time.Duration, tokenType string, role string, scope ...string) (string, error) {
	now := time.Now()
	claims := &TokenClaim{
		id,
		tokenType,
		role,
		scope,
		jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(expiresIn).Unix(),
		}}
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.

	t, err := token.SignedString([]byte(secretJWT))

	if tokenType == "refresh_token" {
		if _, err := redis.SetRefreshToken(id, t, expiresIn); err != nil {
			return "", err
		}
	}

	return t, err
}
