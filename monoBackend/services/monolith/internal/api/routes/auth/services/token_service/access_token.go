package token_service

import (
	"github.com/golang-jwt/jwt/v4"
	"said-and-dot-backend/pkg/config"
	"time"
)

type AccessJwtToken struct {
	jwtToken string
}

func (t *AccessJwtToken) ToString() string {
	return t.jwtToken
}

func NewAccessToken(claims jwt.MapClaims) (*AccessJwtToken, error) {
	generatedTokenString, err := generateJwtToken(claims, getAccessTokenValidityDuration(), getAccessTokenSecretKey())
	if err != nil {
		return nil, err
	}

	accessToken := new(AccessJwtToken)
	accessToken.jwtToken = generatedTokenString

	return accessToken, nil
}

func VerifyAccessToken(tokenString string) (jwt.MapClaims, error) {
	token, tokenClaims, err := ParseJwtToken(tokenString, getAccessTokenSecretKey())
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenSignatureInvalid
	}
	if !tokenClaims.VerifyExpiresAt(time.Now().Unix(), true) {
		return nil, jwt.ErrTokenExpired
	}

	return tokenClaims, nil
}

func getAccessTokenValidityDuration() time.Duration {
	return config.GetDuration("ACCESS_TOKEN_DURATION", time.Duration(time.Now().Add(time.Minute*15).Unix()))
}

func getAccessTokenSecretKey() string {
	return config.GetString("ACCESS_TOKEN_SECRET", "")
}
