package internal

import (
	"context"
	"encoding/json"
	menuhd "food-story/menu-service/internal/adapter/http"
	"food-story/menu-service/internal/adapter/repository"
	"food-story/menu-service/internal/usecase"
	"food-story/pkg/common"
	"food-story/pkg/middleware"
	"food-story/shared/config"
	database "food-story/shared/database/sqlc"
	"food-story/shared/snowflakeid"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"

	dbcfg "food-story/shared/config"
)

const EnvFile = ".env"

type FiberServer struct {
	App *fiber.App

	db *pgxpool.Pool
}

func New() *FiberServer {
	configApp := config.InitConfig(EnvFile)
	app := fiber.New(fiber.Config{
		ServerHeader:             "menu-service",
		AppName:                  "menu-service",
		ErrorHandler:             middleware.HandleError,
		EnableSplittingOnParsers: true,
		JSONEncoder:              json.Marshal,
		JSONDecoder:              json.Unmarshal,
	})

	// add rate limit
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		LimitReached: func(_ *fiber.Ctx) error {
			return fiber.NewError(fiber.StatusTooManyRequests, "Too Many Requests")
		},
	}))

	// add custom CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, Connection",
		AllowMethods: "GET, PUT, POST, PATCH, DELETE, OPTIONS",
	}))

	// add log handler
	app.Use(middleware.LogHandler())

	// connect to database
	configDB := dbcfg.InitDBConfig(EnvFile)
	dbConn, err := configDB.ConnectToDatabase()
	if err != nil {
		panic(err)
	}
	store := database.NewStore(dbConn)

	// Create a new Node with a Node number of 1
	node := snowflakeid.CreateSnowflakeNode(1)
	snowflakeNode := snowflakeid.NewSnowflake(node)

	// init validator
	validator := middleware.NewCustomValidator()

	// init router
	apiV1 := app.Group(common.BasePath)

	// add healthcheck
	apiV1.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(_ *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: common.LivenessEndpoint,
		ReadinessProbe: func(c *fiber.Ctx) bool {
			return readinessDatabase(c.Context(), dbConn)
		},
		ReadinessEndpoint: common.ReadinessEndpoint,
	}))

	registerHandlers(apiV1, store, validator, snowflakeNode, configApp)

	return &FiberServer{
		App: app,
		db:  dbConn,
	}
}

func readinessDatabase(ctx context.Context, dbConn *pgxpool.Pool) bool {
	return dbConn.Ping(ctx) == nil
}

func registerHandlers(router fiber.Router, store database.Store, validator *middleware.CustomValidator, snowflakeNode *snowflakeid.SnowflakeImpl, configApp config.Config) {
	menuRepo := repository.NewRepo(configApp, store, snowflakeNode)
	menuUsecase := usecase.NewUsecase(configApp, *menuRepo)
	menuhd.NewHTTPHandler(router, menuUsecase, validator)
}

func (s *FiberServer) CloseDB() {
	s.db.Close()
}
