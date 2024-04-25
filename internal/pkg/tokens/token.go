package tokens

import (
	"errors"
	"time"

	"github.com/noolingo/user-service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type JWTToken struct {
	auth *domain.AppAuth
}

func New(a *domain.AppAuth) *JWTToken {
	return &JWTToken{auth: a}
}

func (t *JWTToken) NewToken(userID string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"ttl":    time.Now().Add(ttl),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(t.auth.AccessSecretKey))
}

func (t *JWTToken) ParseToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return t.auth.AccessSecretKey, nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("token is invalid") //переписать
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", err
	}

	userID, ok := claims["userID"].(string)
	if !ok {
		return "", errors.New("invalid user id")
	}
	return userID, err
}
