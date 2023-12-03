package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const AccessTokenExpireDuration = time.Hour * 1
const RefreshTokenExpireDuration = time.Hour * 24 * 7

var keyFunc = func(token *jwt.Token) (interface{}, error) {
	return evergreenSecretKey, nil
}
var evergreenSecretKey = []byte("github.com/LeeSeika/evergreen")

type EvergreenClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenToken(userID int64, username string) (string, string, error) {
	accessClaims := EvergreenClaims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(AccessTokenExpireDuration).Unix(),
			Issuer:    "evergreen",
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(evergreenSecretKey)
	if err != nil {
		return "", "", err
	}

	// 因为RefreshToken不需要携带用户信息，索引用jwt.StandardClaims即可
	refreshClaims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(RefreshTokenExpireDuration).Unix(),
		Issuer:    "evergreen",
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(evergreenSecretKey)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func RefreshToken(accessToken, refreshToken string) (string, string, error) {
	// refreshToken无效，直接返回
	if _, err := jwt.Parse(refreshToken, keyFunc); err != nil {
		return "", "", err
	}

	var claims = &EvergreenClaims{}
	_, err := jwt.ParseWithClaims(accessToken, claims, keyFunc)
	v, _ := err.(*jwt.ValidationError)

	if v.Errors == jwt.ValidationErrorExpired {
		return GenToken(claims.UserID, claims.Username)
	}
	return "", "", err
}

func ParseToken(tokenStr string) (*EvergreenClaims, error) {
	var claims = &EvergreenClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, keyFunc)

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
