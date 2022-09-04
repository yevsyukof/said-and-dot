package refresh_token_service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"said-and-dot-backend/pkg/store/postgres"
	"said-and-dot-backend/pkg/validator"
	token_service2 "said-and-dot-backend/services/monolith/internal/api/routes/auth/services/token_service"
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
	db postgres.Store
}

func NewRefreshService(db postgres.Store) RefreshService {
	return refreshService{db: db}
}

// TODO нужно разбивать метод
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

	refreshTokenClaims, err := token_service2.VerifyRefreshToken(refreshInput.RefreshToken)
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

	userID, err := uuid.Parse(refreshTokenClaims["userID"].(string))
	if err != nil {
		return err
	}
	newAccessToken, newRefreshToken, err := token_service2.GenerateNewTokensPair(
		jwt.MapClaims{"userID": userID},
		jwt.MapClaims{"userID": userID})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err,
		})
	}

	if _, err := transaction.Exec("INSERT INTO Refresh_tokens VALUES ($1, $2)",
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
