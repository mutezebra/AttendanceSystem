package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"time"
)

var (
	expireTime = 2 * 24 * time.Hour
	jwtSecret  = []byte("jwt-secret")
)

type Claims struct {
	UID int64
	jwt.StandardClaims
}

func GenerateToken(uid int64) (string, error) {
	expire := time.Now().Add(expireTime).Unix()
	claim := &Claims{
		UID: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire,
			Issuer:    "mutezebra",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	s, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("failed signed token"))
	}
	return s, nil
}

func CheckToken(token string) (int64, bool, error) {
	tokenClaim, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return 0, false, errors.Wrap(err, "failed parse token")
	}
	claim, ok := &Claims{}, false
	if claim, ok = tokenClaim.Claims.(*Claims); !ok {
		return 0, false, errors.Wrap(fmt.Errorf("jurge claim failed"), "")
	}
	return claim.UID, claim.Valid() == nil, nil
}
