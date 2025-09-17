package internal

import (
	"context"
	"encoding/json"
	paymenthd "food-story/payment-service/internal/adapter/http"
	"food-story/payment-service/internal/adapter/repository"
	"food-story/payment-service/internal/usecase"
	"food-story/pkg/common"
	"food-story/pkg/middleware"
	"food-story/shared/config"
	database "food-story/shared/database/sqlc"
	"food-story/shared/snowflakeid"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/swagger"
	"github.com/jackc/pgx/v5/pgxpool"
)

const EnvFile = ".env"
const ServiceName = "payment-service"

type FiberServer struct {
	App    *fiber.App
	Config config.Config

	db *pgxpool.Pool
}

func (s *FiberServer) CloseAllConnection() {
	if s.db != nil {
		s.db.Close()
		log.Println("Database closed")
	}
}

func New() *FiberServer {
	configApp := config.InitConfig(EnvFile)
	configApp.BaseURL = common.BasePath + "/payments"
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

	// add log handler
	app.Use(middleware.LogHandler(configApp.BaseURL))

	// connect to database
	configDB := config.InitDBConfig(EnvFile)
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
			return readiness(c.Context(), dbConn)
		},
		ReadinessEndpoint: common.ReadinessEndpoint,
	}))

	registerHandlers(apiV1, store, validator, snowflakeNode, configApp)

	return &FiberServer{
		App:    app,
		Config: configApp,
		db:     dbConn,
	}
}

func readiness(ctx context.Context, dbConn *pgxpool.Pool) bool {
	return dbConn.Ping(ctx) == nil
}

func registerHandlers(router fiber.Router, store database.Store, validator *middleware.CustomValidator, snowflakeNode *snowflakeid.SnowflakeImpl, configApp config.Config) {
	paymentRepo := repository.NewRepository(configApp, store, snowflakeNode)
	paymentCase := usecase.NewUsecase(configApp, *paymentRepo)

	authInstance := middleware.NewAuthInstance(configApp.KeyCloakCertURL)
	paymenthd.NewHTTPHandler(router, paymentCase, validator, authInstance)
}
