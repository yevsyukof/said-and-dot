package auth

import (
	"github.com/gofiber/fiber/v2"
	"said-and-dot-backend/pkg/store/postgres"
	"said-and-dot-backend/services/monolith/internal/api/routes/auth/services/login_service"
	"said-and-dot-backend/services/monolith/internal/api/routes/auth/services/logout_service"
	"said-and-dot-backend/services/monolith/internal/api/routes/auth/services/refresh_token_service"
	"said-and-dot-backend/services/monolith/internal/api/routes/auth/services/signup_service"
)

type authMiddleware struct {
	loginService   login_service.LoginService
	logoutService  logout_service.LogoutService
	refreshService refresh_token_service.RefreshService
	signupService  signup_service.SignupService
}

func NewAuthMiddleware(db postgres.Store) *authMiddleware {
	return &authMiddleware{
		loginService:   login_service.NewLoginService(db),
		logoutService:  logout_service.NewLogoutService(db),
		refreshService: refresh_token_service.NewRefreshService(db),
		signupService:  signup_service.NewSignupService(db),
	}
}

func SetRoutes(r fiber.Router, db postgres.Store) {
	authMiddleware := NewAuthMiddleware(db)

	r.Post("/login", authMiddleware.loginService.Login)
	r.Post("/logout", authMiddleware.logoutService.Logout)
	r.Post("/refresh-token", authMiddleware.refreshService.Refresh)
	r.Post("/signup", authMiddleware.signupService.Signup)
}
