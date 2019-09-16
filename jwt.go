package main

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2/bson"
)

type tokenClaim struct {
	ID    string   `json:"id"`
	Type  string   `json:"type"`
	Role  string   `json:"role"`
	Scope []string `json:"scope"`
	jwt.StandardClaims
}

func JWTMiddleware() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &tokenClaim{},
		SigningKey: []byte("secret"),
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

func NewToken(id string, expiresIn time.Duration, tokenType string, role string, scope ...string) (string, error) {
	now := time.Now()
	claims := &tokenClaim{
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

	t, err := token.SignedString([]byte("secret"))

	return t, err
}

func ValidateRefreshToken(tokenString string) (bool, *tokenClaim) {
	// validate
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})

	// check
	if err != nil {
		// return false, errors.New("invalid auth token")
		return false, nil
	}
	// check
	if claims, ok := token.Claims.(*tokenClaim); ok && token.Valid {
		fmt.Print(claims.Role)
		return true, claims
	}
	return false, nil
}
