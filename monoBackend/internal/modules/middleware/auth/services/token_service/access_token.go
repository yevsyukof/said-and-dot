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

	if !token.Valid || !tokenClaims.VerifyExpiresAt(time.Now().Unix(), true) {
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

//func GetAccessTokenClaims(accessTokenString string) (jwt.MapClaims, error) {
//	return getTokenClaims(accessTokenString, getAccessTokenSecretKey())
//}
