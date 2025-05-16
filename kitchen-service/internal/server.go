package internal

import (
	"context"
	"encoding/json"
	kitchenhd "food-story/kitchen-service/internal/adapter/http"
	"food-story/kitchen-service/internal/adapter/repository"
	websockethub "food-story/kitchen-service/internal/adapter/websocket"
	"food-story/kitchen-service/internal/usecase"
	"food-story/pkg/common"
	"food-story/pkg/middleware"
	"food-story/shared/config"
	database "food-story/shared/database/sqlc"
	"food-story/shared/kafka"
	"food-story/shared/snowflakeid"
	"github.com/IBM/sarama"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
	"time"

	dbcfg "food-story/shared/config"
)

const EnvFile = ".env"

type FiberServer struct {
	App          *fiber.App
	db           *pgxpool.Pool
	WebsocketHub *websockethub.Hub
	KafkaClient  sarama.ConsumerGroup
}

func New() *FiberServer {
	configApp := config.InitConfig(EnvFile)
	app := fiber.New(fiber.Config{
		ServerHeader:             "kitchen-service",
		AppName:                  "kitchen-service",
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

	// สร้าง WebSocket Hub และเริ่มให้ทำงาน
	hub := websockethub.NewHub()

	// init kafka
	kafkaClient := initKafka(configApp)

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

	registerHandlers(apiV1, store, validator, snowflakeNode, configApp, hub)

	return &FiberServer{
		App:          app,
		db:           dbConn,
		WebsocketHub: hub,
		KafkaClient:  kafkaClient,
	}
}

func readinessDatabase(ctx context.Context, dbConn *pgxpool.Pool) bool {
	return dbConn.Ping(ctx) == nil
}

func registerHandlers(router fiber.Router, store database.Store, validator *middleware.CustomValidator, snowflakeNode *snowflakeid.SnowflakeImpl, configApp config.Config, hub *websockethub.Hub) {

	kitchenRepo := repository.NewRepo(configApp, store, snowflakeNode)
	kitchenUseCase := usecase.NewUsecase(configApp, *kitchenRepo)
	kitchenhd.NewHTTPHandler(router, kitchenUseCase, validator, configApp)

	websockethub.NewWSHandler(router, configApp, hub)
}

func (s *FiberServer) CloseDB() {
	s.db.Close()
}

func (s *FiberServer) CloseWebsocketHub() {
	s.WebsocketHub.Shutdown()
}

func initKafka(configApp config.Config) sarama.ConsumerGroup {
	brokers := strings.Split(configApp.KafkaBrokers, ",")
	client := kafka.InitConsumer(kafka.Group, brokers)
	return client
}
