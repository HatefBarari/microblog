package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"uid"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateTokens(uid, role, accessSecret, refreshSecret string, accMin, refHour int) (acc, ref string, err error) {
	now := time.Now()
	accClaims := Claims{
		UserID: uid,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(accMin) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	accToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accClaims)
	acc, err = accToken.SignedString([]byte(accessSecret))
	if err != nil {
		return "", "", err
	}
	refClaims := Claims{
		UserID: uid,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(refHour) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	refToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refClaims)
	ref, err = refToken.SignedString([]byte(refreshSecret))
	return acc, ref, err
}

func ValidateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}