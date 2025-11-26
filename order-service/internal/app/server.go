package app

import (
	"context"
	"encoding/json"
	"fmt"
	_ "food-story/order-service/docs"
	"food-story/order-service/internal/adapter/cache"
	orderhd "food-story/order-service/internal/adapter/http"
	"food-story/order-service/internal/adapter/queue/producer"
	"food-story/order-service/internal/adapter/repository"
	"food-story/order-service/internal/usecase"
	"food-story/pkg/common"
	"food-story/pkg/middleware"
	"food-story/shared/config"
	database "food-story/shared/database/sqlc"
	"food-story/shared/kafka"
	"food-story/shared/redis"
	"food-story/shared/snowflakeid"
	"log"
	"log/slog"
	"strings"

	"github.com/IBM/sarama"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/swagger"
	"github.com/jackc/pgx/v5/pgxpool"
)

const EnvFile = ".env"
const ServiceName = "order-service"

type FiberServer struct {
	App    *fiber.App
	Config config.Config

	db            *pgxpool.Pool
	redis         *redis.RedisClient
	kafkaProducer sarama.SyncProducer
	clientKafka   sarama.Client
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

	if s.kafkaProducer != nil {
		err := s.kafkaProducer.Close()
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Println("Kafka Producer closed")
	}

	if s.clientKafka != nil {
		err := s.clientKafka.Close()
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Println("Kafka Client closed")
	}

}

func New() (*FiberServer, error) {
	configApp, err := config.InitConfig(EnvFile)
	if err != nil {
		return nil, fmt.Errorf("failed to init config: %w", err)
	}
	configApp.BaseURL = common.BasePath + "/orders"
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

	// init auth
	authInstance, err := middleware.NewAuthInstance(configApp.KeyCloakCertURL)
	if err != nil {
		return nil, fmt.Errorf("failed to init auth instance: %w", err)
	}

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

	// connect to kafka
	brokers := strings.Split(configApp.KafkaBrokers, ",")
	producerKafka, clientKafka, err := kafka.InitProducer(brokers)
	if err != nil {
		return nil, fmt.Errorf("failed to init kafka producer: %w", err)
	}

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
			return readiness(c.Context(), dbConn, redisConn, clientKafka)
		},
		ReadinessEndpoint: common.ReadinessEndpoint,
	}))

	registerHandlers(apiV1, store, validator, snowflakeNode, configApp, redisConn, producerKafka, authInstance)
	return &FiberServer{
		App:           app,
		Config:        configApp,
		db:            dbConn,
		redis:         redisConn,
		kafkaProducer: producerKafka,
		clientKafka:   clientKafka,
	}, nil
}

func readiness(ctx context.Context, dbConn *pgxpool.Pool, redisConn *redis.RedisClient, clientKafka sarama.Client) bool {
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

	_, kafkaErr := clientKafka.Topics()
	if kafkaErr != nil {
		slog.Error("ping kafka", "error: ", kafkaErr)
		return false
	}

	return true
}

func registerHandlers(router fiber.Router, store database.Store, validator *middleware.CustomValidator, snowflakeNode *snowflakeid.SnowflakeImpl, configApp config.Config, redisConn *redis.RedisClient, producerKafka sarama.SyncProducer, authInstance *middleware.AuthInstance) {
	orderQueue := producer.NewQueue(producerKafka)
	orderCache := cache.NewRedisTableCache(redisConn)
	orderRepo := repository.NewRepository(configApp, store, snowflakeNode)
	orderUseCase := usecase.NewUsecase(configApp, *orderRepo, orderCache, orderQueue)

	orderhd.NewHTTPHandler(router, orderUseCase, validator, configApp, authInstance)
}
