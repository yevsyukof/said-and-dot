package token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"time"
)

func VerifyJwtToken(tokenString string, secretKey string) (*jwt.Token, jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, nil, err
	}

	return parsedToken, parsedToken.Claims.(jwt.MapClaims), nil
}

func generateJwtToken(claims jwt.MapClaims, expires time.Duration, secretKey string) (string, error) {
	id, err := gonanoid.New()
	if err != nil {
		return "", err
	}

	claims["id"] = id
	claims["iat"] = time.Now().Unix()
	claims["exp"] = expires

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
