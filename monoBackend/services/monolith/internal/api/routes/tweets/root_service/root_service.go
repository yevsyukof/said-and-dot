package root_service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"said-and-dot-backend/pkg/store/postgres"
	"said-and-dot-backend/pkg/validator"
	"said-and-dot-backend/services/monolith/internal/api/db_middleware"
	"said-and-dot-backend/services/monolith/internal/api/entities"
	"said-and-dot-backend/services/monolith/internal/api/routes/auth/services/token_service"
	"time"
)

type RootService interface {
	Get(ctx *fiber.Ctx) error
	Post(ctx *fiber.Ctx) error
}

type rootService struct {
	db postgres.Store
}

func NewRootService(db postgres.Store) RootService {
	return rootService{db: db}
}

func (rs rootService) Get(ctx *fiber.Ctx) error {
	allTweets, err := db_middleware.GetAllTweets(rs.db)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"allTweets": allTweets,
	})
}

type RootInput struct {
	Tweet string `json:"tweet" validate:"required"`
}

func (mi RootInput) Validate() []*validator.ValidationError {
	return validator.ValidateStruct(mi)
}

func (rs rootService) Post(ctx *fiber.Ctx) error {
	authorizationHeader, contains := ctx.GetReqHeaders()[fiber.HeaderAuthorization]
	if !contains {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "There is no Authorization header in the request",
		})
	}

	var rootInput RootInput
	if err := ctx.BodyParser(&rootInput); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err,
		})
	}
	if validationErrors := rootInput.Validate(); validationErrors != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(validationErrors)
	}

	accessTokenClaims, err := token_service.VerifyAccessToken(authorizationHeader)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid access token",
		})
	}

	userID, err := uuid.Parse(accessTokenClaims["userID"].(string))
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid access token",
		})
	}

	newTweet := entities.Tweet{
		ID:      uuid.UUID{},
		UserID:  userID,
		Tweet:   rootInput.Tweet,
		Created: time.Now(),
		Likes:   nil,
	}

	if err := db_middleware.SaveNewTweet(&newTweet, rs.db); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Tweet successfully saved",
	})
}
