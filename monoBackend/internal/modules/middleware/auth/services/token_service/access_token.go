package token_service

import (
	"github.com/golang-jwt/jwt/v4"
	"said-and-dot-backend/internal/common/config"
	"time"
)

func NewAccessToken(claims jwt.MapClaims) (*JwtToken, error) {
	generatedTokenString, err := generateJwtToken(claims, getAccessTokenValidityDuration(), getAccessTokenSecretKey())
	if err != nil {
		return nil, err
	}
	accessToken := new(JwtToken)
	accessToken.tokenString = generatedTokenString

	return accessToken, nil
}

func VerifyAccessToken(tokenString string) (jwt.MapClaims, error) {
	token, tokenClaims, err := ParseJwtToken(tokenString, getAccessTokenSecretKey())
	if err != nil {
		return nil, err
	}

	// Проверяем есть ли у токена Claims?
	// Ниже - type assertion
	//if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
	//	return nil, jwt.ErrInvalidKey
	//}

	if !token.Valid {
		return nil, jwt.ErrInvalidKey
	}
	return tokenClaims, nil
}

func getAccessTokenValidityDuration() time.Duration {
	return config.GetDuration("ACCESS_TOKEN_DURATION", time.Duration(time.Now().Add(time.Minute*15).Unix()))
}

func getAccessTokenSecretKey() string {
	return config.GetString("ACCESS_TOKEN_SECRET", "")
}
