package token

import (
	"github.com/golang-jwt/jwt/v4"
	"said-and-dot-backend/internal/common/config"
	"time"
)

type AccessToken struct {
	token     string
	expiresAt time.Time
}

func NewAccessToken(claims jwt.MapClaims) (*AccessToken, error) {
	expires := config.GetDuration("ACCESS_TOKEN_DURATION",
		time.Duration(time.Now().Add(time.Minute*15).Unix()))
	secret := config.GetString("ACCESS_TOKEN_SECRET", "")

	token, err := generateJwtToken(claims, expires, secret)
	if err != nil {
		return nil, err
	}

	at := new(AccessToken)
	at.token = token
	at.expiresAt = time.Now().Add(config.GetDuration("ACCESS_TOKEN_DURATION",
		time.Duration(time.Now().Add(time.Minute*15).Unix())))

	return at, nil
}

func VerifyAccessToken(tokenString string) (jwt.MapClaims, error) {
	token, claims, err := VerifyJwtToken(tokenString, config.GetString("ACCESS_TOKEN_SECRET", ""))
	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	return claims, nil
}

func (t *AccessToken) String() string {
	return t.token
}

func (t *AccessToken) ExpiresAt() time.Time {
	return t.expiresAt
}
