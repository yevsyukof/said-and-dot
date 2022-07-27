package token_service

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"time"
)

type JwtToken struct {
	tokenString string
}

func (t *JwtToken) ToString() string {
	return t.tokenString
}

func ParseJwtToken(tokenString string, secretKey string) (*jwt.Token, jwt.MapClaims, error) {
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

func generateJwtToken(claims jwt.MapClaims, validityDuration time.Duration, secretKey string) (string, error) {
	tokenId, err := gonanoid.New()
	if err != nil {
		return "", err
	}

	claims["id"] = tokenId
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(validityDuration).Unix()

	tokenStruct := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	generatedToken, err := tokenStruct.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return generatedToken, nil
}
