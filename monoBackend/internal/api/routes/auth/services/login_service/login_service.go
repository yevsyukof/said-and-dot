package login_service

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"said-and-dot-backend/internal/api/routes/auth/services/token_service"
	"said-and-dot-backend/internal/api/routes/errors"
	"said-and-dot-backend/internal/common/bcrypt"
	"said-and-dot-backend/internal/common/validator"
	"said-and-dot-backend/internal/database"
)

type LoginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (li LoginInput) Validate() []*validator.ValidationError {
	return validator.ValidateStruct(li)
}

type LoginService interface {
	Login(ctx *fiber.Ctx) error
}

type loginService struct {
	db database.Database
}

func NewLoginService(db database.Database) LoginService {
	return loginService{db: db}
}

func (ls loginService) createTokensPair(input LoginInput) (
	*token_service.AccessJwtToken, *token_service.RefreshJwtToken, error) {

	var userId uuid.UUID
	var username, passwordHash string

	if err := ls.db.QueryRow(
		"SELECT id, username, password_hash FROM Users WHERE username = $1",
		input.Username).Scan(&userId, &username, &passwordHash); err != nil {
		return nil, nil, api_errors.ErrUserDoesNotExist
	}

	if !bcrypt.Compare(passwordHash, input.Password) {
		return nil, nil, errors.New("invalid password provided")
	}

	accessToken, err := token_service.NewAccessToken(jwt.MapClaims{
		"userID": userId,
	})
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := token_service.NewRefreshToken(jwt.MapClaims{
		"userID": userId,
	})
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
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
		case errors.Is(err, api_errors.ErrUserDoesNotExist) || errors.Is(err, api_errors.ErrInvalidPassword):
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid email/password",
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "There was a problem on our side",
			})
		}
	}

	if _, err := ls.db.Exec("INSERT INTO Refresh_tokens (user_id, token) VALUES ($1, $2)",
		refreshToken.GetUserID(), refreshToken.ToString()); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "There was a problem on our side",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"accessToken":  accessToken.ToString(),
		"refreshToken": refreshToken.ToString(),
	})
}