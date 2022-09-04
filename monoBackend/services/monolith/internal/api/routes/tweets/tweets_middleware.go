package tweets

import (
	"github.com/gofiber/fiber/v2"
	"said-and-dot-backend/pkg/store/postgres"
	"said-and-dot-backend/services/monolith/internal/api/routes/tweets/my_service"
	"said-and-dot-backend/services/monolith/internal/api/routes/tweets/root_service"
)

type tweetsMiddleware struct {
	rootService root_service.RootService
	myService   my_service.MyService
}

func NewTweetsMiddleware(db postgres.Store) *tweetsMiddleware {
	return &tweetsMiddleware{
		rootService: root_service.NewRootService(db),
		myService:   my_service.NewMyService(db),
	}
}

func SetRoutes(r fiber.Router, db postgres.Store) {
	usersMiddleware := NewTweetsMiddleware(db)

	r.Get("/", usersMiddleware.rootService.Get)
	r.Post("/", usersMiddleware.rootService.Post)
	r.Get("/my", usersMiddleware.myService.Get)
}
