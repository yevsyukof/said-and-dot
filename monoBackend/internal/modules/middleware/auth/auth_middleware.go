package auth

import (
	"github.com/gofiber/fiber/v2"
	"said-and-dot-backend/internal/database"
	"said-and-dot-backend/internal/modules/middleware"
	"said-and-dot-backend/internal/modules/middleware/auth/services/token"
	"strings"
)

type authMiddleware struct{}

func NewAuthMiddleware() middleware.Middleware {
	return authMiddleware{}
}

func (m authMiddleware) Execute() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var accessTokenString string

		if ctx.Cookies("access_token") != "" {
			accessTokenString = ctx.Cookies("access_token")
		} else {
			authorization := ctx.Get("Authorization")

			if len(strings.Split(authorization, " ")) < 2 {
				return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Please login to continue",
				})
			}
			accessTokenString = strings.Split(authorization, " ")[1]
		}

		claims, err := token.VerifyAccessToken(accessTokenString)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Please login to continue",
			})
		}

		userID, ok := claims["userID"]
		if !ok {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Please try and login again",
			})
		}

		email, ok := claims["email"]
		if !ok {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Please try and login again",
			})
		}

		ctx.Locals("userID", userID)
		ctx.Locals("email", email)
		ctx.Locals("accessToken", accessTokenString)

		return ctx.Next()
	}
}

func InitRoutes(r fiber.Router, db database.Database) {
	authMiddleware := NewAuthMiddleware()

	r.Post("/login", buildLoginHandler(db))
	r.Get("/me", authMiddleware.Execute(), buildMeHandler(db)) // аргументы образуют стек функций для роута
	r.Get("/token", buildTokenHandler(db, cache))
	r.Post("/logout", buildLogoutHandler(cache))
}

func buildLoginHandler(db database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		service := service.NewLoginService(db)
		action := action.NewLoginAction(service)

		return action.Execute(c)
	}
}
