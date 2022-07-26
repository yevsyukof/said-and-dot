package login

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"said-and-dot-backend/internal/common/bcrypt"
	"said-and-dot-backend/internal/common/validator"
	"said-and-dot-backend/internal/database"
	auth_errors "said-and-dot-backend/internal/modules/middleware/auth/errors"
	"said-and-dot-backend/internal/modules/middleware/auth/services/token"
)

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (li LoginInput) Validate() []*validator.ValidationError {
	return validator.ValidateStruct(li)
}

type LoginService interface {
	createTokensPair(input LoginInput) (*token.AccessToken, *token.RefreshToken, error)
}

type loginService struct {
	db database.Database
}

func NewLoginService(db database.Database) LoginService {
	return loginService{db: db}
}

func (ls loginService) createTokensPair(input LoginInput) (
	*token.AccessToken, *token.RefreshToken, error) {

	var userId uuid.UUID
	var email, passwordHash string

	if err := ls.db.QueryRow(
		"SELECT id, email, password_hash FROM Users WHERE email = $1",
		input.Email).Scan(&userId, &email, &passwordHash); err != nil {
		return nil, nil, auth_errors.ErrUserDoesNotExist
	}

	if !bcrypt.Compare(passwordHash, input.Password) {
		return nil, nil, errors.New("Invalid password provided")
	}

	at, err := token.NewAccessToken(jwt.MapClaims{
		"userID": userId,
		"email":  email,
	})
	if err != nil {
		return nil, nil, err
	}

	rt, err := token.NewRefreshToken(jwt.MapClaims{
		"userID": userId,
	})
	if err != nil {
		return nil, nil, err
	}

	return at, rt, nil
}

func (ls loginService) Login(ctx *fiber.Ctx) error {
	var loginInput LoginInput

	if err := ctx.BodyParser(&loginInput); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if validationErrors := loginInput.Validate(); validationErrors != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(validationErrors)
	}

	accessToken, refreshToken, err := ls.createTokensPair(loginInput)
	if err != nil {
		switch {
		case errors.Is(err, auth_errors.ErrUserDoesNotExist) || errors.Is(err, auth_errors.ErrInvalidPassword):
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid email/password",
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "There was a problem on our side",
			})
		}
	}

	//ctx.Cookie(&fiber.Cookie{
	//	Name:     "refresh_token",
	//	Value:    refreshToken.String(),
	//	Expires:  refreshToken.ExpiresAt(),
	//	HTTPOnly: true,
	//	Secure:   config.GetString("APP_ENV", "development") == "production",
	//	Path:     "/",
	//	Domain:   config.GetString("APP_DOMAIN", ""),
	//	SameSite: "None",
	//})

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  accessToken.String(),
		"refresh_token": refreshToken.String(),
	})
}
