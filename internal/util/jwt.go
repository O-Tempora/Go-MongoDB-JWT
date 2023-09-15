package util

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTpayload struct {
	User string `json:"user"`
	jwt.RegisteredClaims
}

func GenerateJWT(guid string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodRS512,
		&JWTpayload{
			guid,
			jwt.RegisteredClaims{
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
				NotBefore: jwt.NewNumericDate(time.Now()),
			},
		},
	)
	_, priv, err := GetKeyPair()
	if err != nil {
		return "", err
	}
	st, err := token.SignedString(priv)
	if err != nil {
		return "", err
	}
	return st, nil
}

func ValidateJWT(token string) (*JWTpayload, error) {
	pub, _, err := GetKeyPair()
	if err != nil {
		return nil, err
	}
	t, err := jwt.ParseWithClaims(token, &JWTpayload{}, func(t *jwt.Token) (interface{}, error) {
		return pub, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := t.Claims.(*JWTpayload); ok && t.Valid && claims.ExpiresAt.After(time.Now()) {
		return claims, nil
	}
	return nil, errors.New("invalid jwt")
}

func GenerateRefresh() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	refresh := base64.StdEncoding.EncodeToString(bytes)
	return refresh, nil
}

func GetTokenPair(guid string) (access string, refresh string, err error) {
	access, err = GenerateJWT(guid)
	if err != nil {
		return "", "", err
	}
	refresh, err = GenerateRefresh()
	return access, refresh, err
}

func GetGUIDFromToken(accessToken string) (string, error) {
	token, err := ValidateJWT(accessToken)
	if err != nil {
		return "", err
	}
	return token.User, nil
}
