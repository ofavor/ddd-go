package util

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtClaims struct {
	jwt.RegisteredClaims
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
	Type     int32  `json:"type"`
}

var JwtKey = "historage-base-product"

func JwtEncode(
	uid int64,
	uname string,
	typee int32,
	duration time.Duration,
) (string, error) {
	// Create the Claims
	claims := &JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			Issuer:    "historage",
			Audience:  jwt.ClaimStrings{"*"},
			Subject:   "auth",
			ID:        uuid.NewString(),
		},
		UserId:   uid,
		Username: uname,
		Type:     typee,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JwtKey))
}

func JwtDecode(tokenString string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	if ret, ok := token.Claims.(*JwtClaims); !ok {
		return nil, jwt.ErrTokenInvalidClaims
	} else {
		return ret, nil
	}
}

func JwtValid(jwt *JwtClaims) bool {
	return jwt.ExpiresAt.Unix() > time.Now().Unix()
}
