package app

import (
	"context"
	"encoding/json"
	"fmt"
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
	"log"
	"log/slog"
	"strings"

	"github.com/IBM/sarama"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/swagger"
	"github.com/jackc/pgx/v5/pgxpool"
)

const EnvFile = ".env"
const ServiceName = "kitchen-service"

func (s *FiberServer) CloseAllConnection() {
	if s.db != nil {
		s.db.Close()
		log.Println("Database closed")
	}

	if s.WebsocketHub != nil {
		s.WebsocketHub.Shutdown()
		log.Println("WebsocketHub closed")
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

type FiberServer struct {
	App           *fiber.App
	Config        config.Config
	KafkaConsumer sarama.ConsumerGroup
	WebsocketHub  *websockethub.Hub

	db          *pgxpool.Pool
	clientKafka sarama.Client
}

func New() (*FiberServer, error) {
	configApp, err := config.InitConfig(EnvFile)
	if err != nil {
		return nil, fmt.Errorf("failed to init config: %w", err)
	}
	configApp.BaseURL = common.BasePath + "/kitchen"
	app := fiber.New(fiber.Config{
		ServerHeader:             ServiceName,
		AppName:                  ServiceName,
		ErrorHandler:             middleware.HandleError,
		EnableSplittingOnParsers: true,
		JSONEncoder:              json.Marshal,
		JSONDecoder:              json.Unmarshal,
	})

	// add custom CORS
	app.Use(cors.New(middleware.DefaultCorsConfig()))

	// add log handler
	app.Use(middleware.LogHandler(configApp.BaseURL))

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

	// สร้าง WebSocket Hub และเริ่มให้ทำงาน
	hub := websockethub.NewHub()

	// init kafka
	brokers := strings.Split(configApp.KafkaBrokers, ",")
	consumerClient, clientKafka, err := kafka.InitConsumer(kafka.Group, brokers)
	if err != nil {
		return nil, fmt.Errorf("failed to init kafka consumer: %w", err)
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
			return readiness(c.Context(), dbConn, clientKafka)
		},
		ReadinessEndpoint: common.ReadinessEndpoint,
	}))

	registerHandlers(apiV1, store, validator, snowflakeNode, configApp, hub, authInstance)
	return &FiberServer{
		App:           app,
		Config:        configApp,
		WebsocketHub:  hub,
		KafkaConsumer: consumerClient,

		db:          dbConn,
		clientKafka: clientKafka,
	}, nil
}

func readiness(ctx context.Context, dbConn *pgxpool.Pool, clientKafka sarama.Client) bool {
	dbErr := dbConn.Ping(ctx)
	if dbErr != nil {
		slog.Error("ping database", "error: ", dbErr)
		return false
	}

	_, kafkaErr := clientKafka.Topics()
	if kafkaErr != nil {
		slog.Error("ping kafka", "error: ", kafkaErr)
		return false
	}

	return true
}

func registerHandlers(router fiber.Router, store database.Store, validator *middleware.CustomValidator, snowflakeNode *snowflakeid.SnowflakeImpl, configApp config.Config, hub *websockethub.Hub, authInstance *middleware.AuthInstance) {

	kitchenRepo := repository.NewRepository(configApp, store, snowflakeNode)
	kitchenUseCase := usecase.NewUsecase(configApp, *kitchenRepo)

	kitchenhd.NewHTTPHandler(router, kitchenUseCase, validator, configApp, authInstance)

	websockethub.NewWSHandler(router, configApp, hub)
}
