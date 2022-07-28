package token_service

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"time"
)

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

func getTokenClaims(tokenString, secretKey string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	return parsedToken.Claims.(jwt.MapClaims), nil
}

func GenerateNewTokensPair(refreshTokenClaims jwt.MapClaims,
	accessTokenClaims jwt.MapClaims) (*AccessJwtToken, *RefreshJwtToken, error) {

	accessToken, err := NewAccessToken(accessTokenClaims)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := NewRefreshToken(refreshTokenClaims)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}
