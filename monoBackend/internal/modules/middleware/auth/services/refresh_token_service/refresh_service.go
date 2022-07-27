package refresh_token_service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"said-and-dot-backend/internal/common/validator"
	"said-and-dot-backend/internal/database"
	"said-and-dot-backend/internal/modules/middleware/auth/services/token_service"
)

type RefreshInput struct {
	RefreshToken string `json:"refreshToken" validate:"required,jwt"`
}

func (ri RefreshInput) Validate() []*validator.ValidationError {
	return validator.ValidateStruct(ri)
}

type RefreshService interface {
	Refresh(ctx *fiber.Ctx) error
}

type refreshService struct {
	db database.Database
}

func NewRefreshService(db database.Database) RefreshService {
	return refreshService{db: db}
}

func (rs refreshService) Refresh(ctx *fiber.Ctx) error {
	var refreshInput RefreshInput

	if err := ctx.BodyParser(&refreshInput); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if validationErrors := refreshInput.Validate(); validationErrors != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(validationErrors)
	}

	refreshTokenClaims, err := token_service.VerifyRefreshToken(refreshInput.RefreshToken)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": err,
		})
	}

	var isTokenExists bool
	if err := rs.db.QueryRow("SELECT EXISTS (SELECT 1 FROM Refresh_tokens WHERE user_id = $1 AND token = $2)",
		refreshTokenClaims["userID"], refreshInput.RefreshToken).Scan(&isTokenExists); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err,
		})
	} else if !isTokenExists {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "This token does not exists",
		})
	}

	transaction, err := rs.db.BeginTx()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err,
		})
	}

	if _, err := transaction.Exec("DELETE FROM Refresh_tokens WHERE user_id = $1 AND token = $2",
		refreshTokenClaims["userID"], refreshInput.RefreshToken); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err,
		})
	}

	newAccessToken, newRefreshToken, err := token_service.GenerateNewTokensPair(
		jwt.MapClaims{"userID": refreshTokenClaims["userID"]}, jwt.MapClaims{"userID": refreshTokenClaims["userID"]})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err,
		})
	}

	if _, err := transaction.Exec("INSERT INTO \"Refresh_tokens\" VALUES ($1, $2)",
		refreshTokenClaims["userID"], newRefreshToken.ToString()); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err,
		})
	}

	if err := transaction.Commit(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":      "Success refresh token",
		"accessToken":  newAccessToken.ToString(),
		"refreshToken": newRefreshToken.ToString(),
	})
}
