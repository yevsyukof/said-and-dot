package token_service

import (
	"github.com/golang-jwt/jwt/v4"
	"said-and-dot-backend/internal/common/config"
	"time"
)

func NewRefreshToken(claims jwt.MapClaims) (*JwtToken, error) {
	generatedTokenString, err := generateJwtToken(claims, getRefreshTokenValidityDuration(), getRefreshTokenSecretKey())
	if err != nil {
		return nil, err
	}

	refreshToken := new(JwtToken)
	refreshToken.tokenString = generatedTokenString

	return refreshToken, nil
}

func VerifyRefreshToken(tokenStr string) (jwt.MapClaims, error) {
	token, tokenClaims, err := ParseJwtToken(tokenStr, getRefreshTokenSecretKey())
	if err != nil {
		return nil, err
	}

	if !token.Valid || !tokenClaims.VerifyExpiresAt(time.Now().Unix(), true) {
		return nil, jwt.ErrInvalidKey
	}
	return tokenClaims, nil
}

func getRefreshTokenValidityDuration() time.Duration {
	return config.GetDuration("REFRESH_TOKEN_DURATION", time.Duration(time.Now().Add(time.Hour*24*7).Unix()))
}

func getRefreshTokenSecretKey() string {
	return config.GetString("REFRESH_TOKEN_SECRET", "")
}

//func GetRefreshTokenClaims(refreshTokenString string) (jwt.MapClaims, error) {
//	return getTokenClaims(refreshTokenString, getRefreshTokenSecretKey())
//}
