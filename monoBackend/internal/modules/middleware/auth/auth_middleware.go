package auth

import (
	"github.com/gofiber/fiber/v2"
	"said-and-dot-backend/internal/database"
	"said-and-dot-backend/internal/modules/middleware/auth/services/login_service"
	"said-and-dot-backend/internal/modules/middleware/auth/services/logout_service"
	"said-and-dot-backend/internal/modules/middleware/auth/services/refresh_token_service"
	"said-and-dot-backend/internal/modules/middleware/auth/services/token_service"
	"strings"
)

type authMiddleware struct {
	loginService   login_service.LoginService
	logoutService  logout_service.LogoutService
	refreshService refresh_token_service.RefreshService
}

func NewAuthMiddleware(db database.Database) *authMiddleware {
	return &authMiddleware{
		loginService:   login_service.NewLoginService(db),
		logoutService:  logout_service.NewLogoutService(db),
		refreshService: refresh_token_service.NewRefreshService(db),
	}
}

// TODO нахуйя он нужон
func (m authMiddleware) Execute() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var accessTokenString string

		if ctx.Cookies("access_token") != "" {
			accessTokenString = ctx.Cookies("access_token")
		} else {
			authorization := ctx.Get("Authorization")

			if len(strings.Split(authorization, " ")) < 2 {
				return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Please login_service to continue",
				})
			}
			accessTokenString = strings.Split(authorization, " ")[1]
		}

		claims, err := token_service.VerifyAccessToken(accessTokenString)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Please login_service to continue",
			})
		}

		userID, ok := claims["userID"]
		if !ok {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Please try and login_service again",
			})
		}

		email, ok := claims["email"]
		if !ok {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Please try and login_service again",
			})
		}

		ctx.Locals("userID", userID)
		ctx.Locals("email", email)
		ctx.Locals("accessToken", accessTokenString)

		return ctx.Next()
	}
}

func SetRoutes(r fiber.Router, db database.Database) {
	authMiddleware := NewAuthMiddleware(db)

	r.Post("/login", authMiddleware.loginService.Login)
	r.Post("/logout", authMiddleware.logoutService.Logout)
	r.Post("/refresh-token", authMiddleware.refreshService.Refresh)

	//r.Post("/signup")
}
