package users

import (
	"github.com/gofiber/fiber/v2"
	"said-and-dot-backend/pkg/store/postgres"
	"said-and-dot-backend/services/monolith/internal/api/routes/users/services/me_service"
)

type usersMiddleware struct {
	meService me_service.MeService
}

func NewUsersMiddleware(db postgres.Store) *usersMiddleware {
	return &usersMiddleware{
		meService: me_service.NewMeService(db),
	}
}

func SetRoutes(r fiber.Router, db postgres.Store) {
	usersMiddleware := NewUsersMiddleware(db)

	r.Get("/me", usersMiddleware.meService.Get)
}
