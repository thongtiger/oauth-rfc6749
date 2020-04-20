package auth

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type TokenClaim struct {
	ID    string   `json:"id"`
	Type  string   `json:"type"`
	Role  string   `json:"role"`
	Scope []string `json:"scope"`
	jwt.StandardClaims
}

type Oauth2 struct {
	Username     string `json:"username,omitempty" form:"username" query:"username"`
	Password     string `json:"password,omitempty" form:"password" query:"password"`
	GrantType    string `json:"grant_type" form:"grant_type" query:"grant_type"`
	RefreshToken string `json:"refresh_token" form:"refresh_token" query:"refresh_token"`
}

type User struct {
	ID             primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Role           string             `json:"role" bson:"role,omitempty"`
	Scope          []string           `json:"scope"`
	Username       string             `json:"username" bson:"username,omitempty"`
	Password       string             `json:"password" bson:"password,omitempty"`
	Name           string             `json:"name" bson:"name,omitempty"`
	CreateTime     time.Time          `json:"createTime" bson:"createTime"`
	LatestLoggedin time.Time          `json:"latestLoggedin" bson:"latestLoggedin"`
}

// BcryptCost : Cost
const BcryptCost = 13

// VerifyPassword : checking
func (u *User) VerifyPassword(input string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// HashingPassword : when set to model
func (u *User) HashingPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), BcryptCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}
