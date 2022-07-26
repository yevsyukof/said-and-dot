package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.uber.org/zap"
	"said-and-dot-backend/internal/common/config"
	"said-and-dot-backend/internal/database"
	"time"
)

type Server struct {
	engine *fiber.App
	config *Config
	db     database.Database
	log    *zap.SugaredLogger
}

// New creates a new instance of a fiber web server
func New(db database.Database, log *zap.SugaredLogger, appCfg *Config, engineCfg ...fiber.Config) *Server {
	appCfg.init()
	return &Server{
		engine: fiber.New(engineCfg...),
		config: appCfg,
		db:     db,
		log:    log,
	}
}

func (s *Server) Listen() {
	s.initMiddlewares()
	s.initRouteGroups()

	if !fiber.IsChild() {
		s.log.Infof("Starting up %s", s.config.AppName)
	}
	if err := s.engine.Listen(fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)); err != nil {
		s.log.Error(err)
	}
}

func (s *Server) initMiddlewares() {
	s.engine.Use(cors.New(cors.Config{
		AllowOrigins:     config.GetString("APP_DOMAIN", "*"),
		AllowCredentials: true,
		AllowHeaders:     "Content-Type",
	}))

	//s.engine.Use(helmet.New()) // TODO

	s.engine.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	s.engine.Use(limiter.New(limiter.Config{
		Max:        60,
		Expiration: 1 * time.Minute,
	}))

	s.engine.Use(logger.New(logger.Config{
		Format: "${green}${time}${reset} | ${status} | ${cyan}${latency}${reset}	-	${host} | ${yellow}${method}${reset} | ${path} ${queryParams}\n",
	}))
}

func (s *Server) initRouteGroups() {
	auth.Routes(
		s.engine.Group("/auth"),
		s.db,
		s.cache)

	//tweet.Routes(
	//	s.engine.Group("/tweets"),
	//	s.db,
	//)
	//
	//user.Routes(
	//	s.engine.Group("/users"),
	//	s.db,
	//)
	//
	//relationship.Routes(
	//	s.engine.Group("/relationships"),
	//	s.db,
	//)
}
