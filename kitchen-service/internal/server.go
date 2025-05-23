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
	"github.com/gofiber/swagger"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"log/slog"
	"strings"
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

func New() *FiberServer {
	configApp := config.InitConfig(EnvFile)
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

	// add log handler
	app.Use(middleware.LogHandler(configApp.BaseURL))

	// connect to database
	configDB := config.InitDBConfig(EnvFile)
	dbConn, err := configDB.ConnectToDatabase()
	if err != nil {
		panic(err)
	}
	store := database.NewStore(dbConn)

	// สร้าง WebSocket Hub และเริ่มให้ทำงาน
	hub := websockethub.NewHub()

	// init kafka
	brokers := strings.Split(configApp.KafkaBrokers, ",")
	consumerClient, clientKafka, err := kafka.InitConsumer(kafka.Group, brokers)
	if err != nil {
		panic(err)
	}

	// Create a new Node with a Node number of 1
	node := snowflakeid.CreateSnowflakeNode(1)
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

	registerHandlers(apiV1, store, validator, snowflakeNode, configApp, hub)
	return &FiberServer{
		App:           app,
		Config:        configApp,
		WebsocketHub:  hub,
		KafkaConsumer: consumerClient,

		db:          dbConn,
		clientKafka: clientKafka,
	}
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

func registerHandlers(router fiber.Router, store database.Store, validator *middleware.CustomValidator, snowflakeNode *snowflakeid.SnowflakeImpl, configApp config.Config, hub *websockethub.Hub) {

	kitchenRepo := repository.NewRepository(configApp, store, snowflakeNode)
	kitchenUseCase := usecase.NewUsecase(configApp, *kitchenRepo)
	kitchenhd.NewHTTPHandler(router, kitchenUseCase, validator, configApp)

	websockethub.NewWSHandler(router, configApp, hub)
}
