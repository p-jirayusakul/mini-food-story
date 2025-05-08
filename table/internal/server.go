package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"food-story/pkg/common"
	"food-story/pkg/middleware"
	dbcfg "food-story/shared/config"
	"food-story/shared/database/sqlc"
	"food-story/shared/snowflakeid"
	tablehd "food-story/table/internal/adapter/http"
	"food-story/table/internal/config"
	"food-story/table/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

const EnvFile = ".env"

type FiberServer struct {
	App *fiber.App

	db *pgxpool.Pool
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
	dbConn, err := connectToDatabase(configDB)
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

func connectToDatabase(configDB dbcfg.DBConfig) (*pgxpool.Pool, error) {
	source := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s search_path=%s sslmode=disable TimeZone=Asia/Bangkok", configDB.DBUsername, configDB.DBPassword, configDB.DBHost, configDB.DBPort, configDB.DBDatabase, configDB.DBSchema)
	dbConn, err := pgxpool.New(context.Background(), source)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	return dbConn, nil
}

func readinessDatabase(ctx context.Context, dbConn *pgxpool.Pool) bool {
	return dbConn.Ping(ctx) == nil
}

func registerHandlers(router fiber.Router, store database.Store, validator *middleware.CustomValidator, snowflakeNode *snowflakeid.SnowflakeImpl, configApp config.Config) {
	tableUseCase := usecase.NewUsecase(configApp, store, snowflakeNode)
	tablehd.NewHTTPHandler(router, tableUseCase, validator)
}

func (s *FiberServer) CloseDB() {
	s.db.Close()
}
