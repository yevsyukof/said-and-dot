package my_service

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"said-and-dot-backend/internal/api/db_middleware"
	"said-and-dot-backend/internal/api/routes/auth/services/token_service"
	api_errors "said-and-dot-backend/internal/api/routes/errors"
	"said-and-dot-backend/internal/common/validator"
	"said-and-dot-backend/internal/database"
)

type MyInput struct {
	AccessToken string `json:"accessToken" validate:"required,jwt"`
}

func (mi MyInput) Validate() []*validator.ValidationError {
	return validator.ValidateStruct(mi)
}

type MyService interface {
	Get(ctx *fiber.Ctx) error
}

type myService struct {
	db database.Database
}

func NewMyService(db database.Database) MyService {
	return myService{db: db}
}

func (ms myService) Get(ctx *fiber.Ctx) error {
	authorizationHeader, contains := ctx.GetReqHeaders()[fiber.HeaderAuthorization]
	if !contains {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "There is no Authorization header in the request",
		})
	}

	meInput := MyInput{AccessToken: authorizationHeader}
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

	userTweets, err := db_middleware.GetUserTweets(user, ms.db)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "There was a problem on our side with database",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{ // TODO
		"userTweets": userTweets,
	})
}
