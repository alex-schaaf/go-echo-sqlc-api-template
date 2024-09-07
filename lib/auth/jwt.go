package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokenCookie(userID int64) (http.Cookie, error) {
	tokenStr, expiresAt, err := CreateToken(userID)
	if err != nil {
		return http.Cookie{}, err
	}
	cookie := http.Cookie{
		Name:     "token",
		Value:    tokenStr,
		Expires:  expiresAt,
		HttpOnly: true,                    // Prevent XSS attacks
		SameSite: http.SameSiteStrictMode, // Prevent CSRF attacks
	}
	return cookie, nil
}

func CreateToken(userID int64) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Hour * 24)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": expiresAt.Unix(),
	})
	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", expiresAt, err
	}
	return tokenStr, expiresAt, nil
}

func IsValidToken(tokenStr string) bool {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return false
	}
	return token.Valid
}

func GetTokenPayload(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims.(jwt.MapClaims), nil
}
