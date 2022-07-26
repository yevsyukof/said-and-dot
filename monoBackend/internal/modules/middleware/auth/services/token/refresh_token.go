package token

import (
	jwtgo "github.com/golang-jwt/jwt/v4"
	"said-and-dot-backend/internal/common/config"
	"time"
)

type RefreshToken struct {
	token     string
	expiresAt time.Time
}

func NewRefreshToken(claims jwtgo.MapClaims) (*RefreshToken, error) {
	exp := config.GetDuration("REFRESH_TOKEN_DURATION",
		time.Duration(time.Now().Add(time.Hour*24*7).Unix()))
	secretKey := config.GetString("REFRESH_TOKEN_SECRET", "")

	token, err := generateJwtToken(claims, exp, secretKey)
	if err != nil {
		return nil, err
	}

	rt := new(RefreshToken)
	rt.token = token
	rt.expiresAt = time.Now().Add(config.GetDuration("REFRESH_TOKEN_DURATION",
		time.Duration(time.Now().Add(time.Hour*24*7).Unix())))

	return rt, nil
}

func VerifyRefreshToken(tokenStr string) (jwtgo.MapClaims, error) {
	token, claims, err := VerifyJwtToken(tokenStr, config.GetString("REFRESH_TOKEN_SECRET", ""))
	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwtgo.Claims); !ok && !token.Valid {
		return nil, jwtgo.ErrInvalidKey
	}

	return claims, nil
}

func (t *RefreshToken) String() string {
	return t.token
}

func (t *RefreshToken) ExpiresAt() time.Time {
	return t.expiresAt
}
