package internal

import (
	"context"
	"encoding/json"
	"food-story/order-service/internal/adapter/cache"
	orderhd "food-story/order-service/internal/adapter/http"
	"food-story/order-service/internal/adapter/queue/producer"
	"food-story/order-service/internal/adapter/repository"
	"food-story/order-service/internal/usecase"
	"food-story/pkg/common"
	"food-story/pkg/middleware"
	"food-story/shared/config"
	dbcfg "food-story/shared/config"
	database "food-story/shared/database/sqlc"
	"food-story/shared/redis"
	"food-story/shared/snowflakeid"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
	"time"
)

const EnvFile = ".env"

type FiberServer struct {
	App *fiber.App

	db    *pgxpool.Pool
	redis *redis.RedisClient
}

func New() *FiberServer {
	configApp := config.InitConfig(EnvFile)
	app := fiber.New(fiber.Config{
		ServerHeader:             "mini-food-story",
		AppName:                  "mini-food-story",
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

	// connect to redis
	redisConn := redis.NewRedisClient(configApp.RedisAddress, configApp.RedisPassword, 0)

	// connect to kafka
	brokers := strings.Split(configApp.KafkaBrokers, ",")
	kafkaConn, err := producer.NewOrderProducer(brokers)
	if err != nil {
		panic(err)
	}

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

	registerHandlers(apiV1, store, validator, snowflakeNode, configApp, redisConn, kafkaConn)

	return &FiberServer{
		App:   app,
		db:    dbConn,
		redis: redisConn,
	}
}

func readinessDatabase(ctx context.Context, dbConn *pgxpool.Pool) bool {
	return dbConn.Ping(ctx) == nil
}

func registerHandlers(router fiber.Router, store database.Store, validator *middleware.CustomValidator, snowflakeNode *snowflakeid.SnowflakeImpl, configApp config.Config, redisConn *redis.RedisClient, kafkaConn *producer.OrderProducer) {
	orderCache := cache.NewRedisTableCache(redisConn)
	orderRepo := repository.NewRepo(configApp, store, snowflakeNode)
	orderUseCase := usecase.NewUsecase(configApp, *orderRepo, orderCache, *kafkaConn)
	orderhd.NewHTTPHandler(router, orderUseCase, validator, configApp)
}

func (s *FiberServer) CloseDB() {
	s.db.Close()
}

func (s *FiberServer) CloseRedis() {
	s.redis.Close()
}
