package app

import (
	"context"
	"encoding/json"
	"fmt"
	"food-story/pkg/common"
	"food-story/pkg/middleware"
	"food-story/shared/config"
	database "food-story/shared/database/sqlc"
	"food-story/shared/redis"
	"food-story/shared/snowflakeid"
	"food-story/table-service/internal/adapter/cache"
	tablehd "food-story/table-service/internal/adapter/http"
	"food-story/table-service/internal/adapter/repository"
	"food-story/table-service/internal/usecase"
	"log"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/swagger"
	"github.com/jackc/pgx/v5/pgxpool"
)

const EnvFile = ".env"
const ServiceName = "table-service"

type FiberServer struct {
	App    *fiber.App
	Config config.Config

	db    *pgxpool.Pool
	redis *redis.RedisClient
}

func (s *FiberServer) CloseAllConnection() {
	if s.db != nil {
		s.db.Close()
		log.Println("Database closed")
	}

	if s.redis != nil {
		s.redis.Close()
		log.Println("Redis closed")
	}
}

func New() (*FiberServer, error) {
	configApp, err := config.InitConfig(EnvFile)
	if err != nil {
		return nil, fmt.Errorf("failed to init config: %w", err)
	}
	configApp.BaseURL = common.BasePath + "/tables"

	app := fiber.New(fiber.Config{
		ServerHeader:             ServiceName,
		AppName:                  ServiceName,
		ErrorHandler:             middleware.HandleError,
		EnableSplittingOnParsers: true,
		JSONEncoder:              json.Marshal,
		JSONDecoder:              json.Unmarshal,
	})

	// add rate limit
	app.Use(limiter.New(middleware.DefaultLimiter()))

	// add custom CORS
	app.Use(cors.New(middleware.DefaultCorsConfig()))

	// add request id
	app.Use(middleware.RequestIDMiddleware())

	// add log handler
	app.Use(middleware.LogHandler(configApp.BaseURL))

	// connect to database
	configDB, err := config.InitDBConfig(EnvFile)
	if err != nil {
		return nil, fmt.Errorf("failed to init db config: %w", err)
	}

	dbConn, err := configDB.ConnectToDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	store := database.NewStore(dbConn)

	// connect to redis
	redisConn := redis.NewRedisClient(configApp.RedisAddress, configApp.RedisPassword, 0)

	// Create a new Node with a Node number of 1
	node, err := snowflakeid.CreateSnowflakeNode(1)
	if err != nil {
		return nil, fmt.Errorf("failed to create snowflake: %w", err)
	}
	snowflakeNode := snowflakeid.NewSnowflake(node)

	// init validator
	validator := middleware.NewCustomValidator()

	// init router
	apiV1 := app.Group(configApp.BaseURL)

	// init swagger endpoint
	apiV1.Get(common.SwaggerEndpoint+"/*", swagger.HandlerDefault)

	// add healthcheck
	apiV1.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(_ *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: common.LivenessEndpoint,
		ReadinessProbe: func(c *fiber.Ctx) bool {
			return readiness(c.Context(), dbConn, redisConn)
		},
		ReadinessEndpoint: common.ReadinessEndpoint,
	}))

	registerHandlers(apiV1, store, validator, snowflakeNode, configApp, redisConn)

	return &FiberServer{
		App:    app,
		Config: configApp,
		db:     dbConn,
		redis:  redisConn,
	}, nil
}

func readiness(ctx context.Context, dbConn *pgxpool.Pool, redisConn *redis.RedisClient) bool {
	dbErr := dbConn.Ping(ctx)
	if dbErr != nil {
		slog.Error("ping database", "error: ", dbErr)
		return false
	}

	redisErr := redisConn.Client.Ping(ctx).Err()
	if redisErr != nil {
		slog.Error("ping redis", "error: ", redisErr)
		return false
	}

	return true
}

func registerHandlers(router fiber.Router, store database.Store, validator *middleware.CustomValidator, snowflakeNode *snowflakeid.SnowflakeImpl, configApp config.Config, redisConn *redis.RedisClient) {
	tableCache := cache.NewRedisTableCache(redisConn)
	tableRepo := repository.NewRepository(configApp, store, snowflakeNode)
	tableUseCase := usecase.NewUseCase(configApp, *tableRepo, tableCache)

	tablehd.NewHTTPHandler(router, tableUseCase, validator)
}
