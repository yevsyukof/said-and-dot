package users

import (
	"github.com/gofiber/fiber/v2"
	"said-and-dot-backend/internal/api/routes/users/services/me_service"
	"said-and-dot-backend/internal/database"
)

type usersMiddleware struct {
	meService me_service.MeService
}

func NewUsersMiddleware(db database.Database) *usersMiddleware {
	return &usersMiddleware{
		meService: me_service.NewMeService(db),
	}
}

func SetRoutes(r fiber.Router, db database.Database) {
	usersMiddleware := NewUsersMiddleware(db)

	r.Get("/me", usersMiddleware.meService.Get)
}
