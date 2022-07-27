package logout_service

import (
	"github.com/gofiber/fiber/v2"
	"said-and-dot-backend/internal/common/validator"
	"said-and-dot-backend/internal/database"
	"said-and-dot-backend/internal/modules/middleware/auth/services/token_service"
)

type LogoutInput struct {
	RefreshToken string `json:"refreshToken" validate:"required,jwt"`
}

func (li LogoutInput) Validate() []*validator.ValidationError {
	return validator.ValidateStruct(li)
}

type LogoutService interface {
	Logout(ctx *fiber.Ctx) error
}

type logoutService struct {
	db database.Database
}

func NewLogoutService(db database.Database) LogoutService {
	return logoutService{db: db}
}

func (ls logoutService) Logout(ctx *fiber.Ctx) error {
	var logoutInput LogoutInput

	if err := ctx.BodyParser(&logoutInput); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if validationErrors := logoutInput.Validate(); validationErrors != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(validationErrors)
	}

	refreshTokenClaims, err := token_service.VerifyRefreshToken(logoutInput.RefreshToken)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": err,
		})
	}

	if _, err := ls.db.Exec("DELETE FROM Refresh_tokens WHERE user_id = $1 AND token = $2",
		refreshTokenClaims["userID"], logoutInput.RefreshToken); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success logout",
	})
}
