package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"dicetales.com/pkg/ctxdata"
	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrInvalidToken = errors.New("invalid jwt token")
)

// GenerateAccessToken generates a short-lived access token
func GenerateAccessToken(secretKey string, iat, seconds int64, uid string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["typ"] = "access"
	claims[ctxdata.Identify] = uid

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(secretKey))
}

// GenerateRefreshToken generates a long-lived refresh token
func GenerateRefreshToken(secretKey string, iat, seconds int64, uid, jti string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["typ"] = "refresh"
	claims["jti"] = jti
	claims[ctxdata.Identify] = uid

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(secretKey))
}

// ParseRefreshToken parses and validates a refresh token
func ParseRefreshToken(tokenStr, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if typ, ok := claims["typ"].(string); !ok || typ != "refresh" {
			return nil, errors.New("not a refresh token")
		}
		return claims, nil
	}
	return nil, ErrInvalidToken
}

// HashToken string to SHA256 string for redis storage
func HashToken(token string) string {
	h := sha256.New()
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}
