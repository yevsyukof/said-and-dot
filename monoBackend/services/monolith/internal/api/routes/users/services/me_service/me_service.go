package me_service

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"said-and-dot-backend/pkg/store/postgres"
	"said-and-dot-backend/pkg/validator"
	"said-and-dot-backend/services/monolith/internal/api/db_middleware"
	"said-and-dot-backend/services/monolith/internal/api/routes/auth/services/token_service"
	"said-and-dot-backend/services/monolith/internal/api/routes/errors"
)

type MeInput struct {
	AccessToken string `json:"accessToken" validate:"required,jwt"`
}

func (mi MeInput) Validate() []*validator.ValidationError {
	return validator.ValidateStruct(mi)
}

type MeService interface {
	Get(ctx *fiber.Ctx) error
}

type meService struct {
	db postgres.Store
}

func NewMeService(db postgres.Store) MeService {
	return meService{db: db}
}

func (ms meService) Get(ctx *fiber.Ctx) error {
	authorizationHeader, contains := ctx.GetReqHeaders()[fiber.HeaderAuthorization]
	if !contains {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "There is no Authorization header in the request",
		})
	}

	meInput := MeInput{AccessToken: authorizationHeader}
	if validationErrors := meInput.Validate(); validationErrors != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(validationErrors)
	}

	accessTokenClaims, err := token_service.VerifyAccessToken(meInput.AccessToken)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid access token",
		})
	}

	user, err := db_middleware.GetUserByID(accessTokenClaims["userID"].(string), ms.db)
	if err != nil {
		switch {
		case errors.Is(err, api_errors.ErrUserDoesNotExist):
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": err,
			})
		case errors.Is(err, api_errors.ErrDatabaseError):
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err,
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "There was a problem on our side",
			})
		}
	}

	userFollowers, err := db_middleware.GetFollowersByUserID(user.ID, ms.db)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "There was a problem on our side with database",
		})
	}

	follows, err := db_middleware.GetFollowsByUserID(user.ID, ms.db)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "There was a problem on our side with database",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{ // TODO
		"userData":  user,
		"followers": userFollowers,
		"follows":   follows,
	})
}
