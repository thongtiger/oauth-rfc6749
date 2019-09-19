package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
)

type tokenClaim struct {
	ID    string   `json:"id"`
	Type  string   `json:"type"`
	Role  string   `json:"role"`
	Scope []string `json:"scope"`
	jwt.StandardClaims
}

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
