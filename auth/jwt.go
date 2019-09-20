package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2/bson"
)

const accessTokenDuration = time.Duration(time.Minute * 15)
const refreshTokenDuration = time.Duration(time.Minute * 30)

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
	if claims, ok := token.Claims.(*TokenClaim); ok && token.Valid {
		fmt.Print(claims.Role)
		return true, claims
	}
	return false, nil
}

func JWTMiddleware() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &TokenClaim{},
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

	t, err := token.SignedString([]byte("secret"))

	if tokenType == "refresh_token" {

	}

	return t, err
}
