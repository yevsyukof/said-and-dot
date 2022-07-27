package logout_service

import (
	"github.com/gofiber/fiber/v2"
	"said-and-dot-backend/internal/common/validator"
	"said-and-dot-backend/internal/database"
)

type LogoutInput struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
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

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success logout",
	})
}
