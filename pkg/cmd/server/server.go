package cmd

import (
	"context"
	"fmt"
	"galasejahtera/pkg/handlers"
	"galasejahtera/pkg/logger"
	model2 "galasejahtera/pkg/model"
	"galasejahtera/pkg/protocol/grpc"
	"galasejahtera/pkg/protocol/rest"
	"galasejahtera/pkg/utility"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

// Config is configuration for Server
type Config struct {
	GRPCPort      string
	HTTPPort      string
	LogLevel      int
	LogTimeFormat string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	// get configuration
	cfg := &Config{GRPCPort: "10001", HTTPPort: "10002", LogLevel: -1, LogTimeFormat: "02 Jan 2006 15:04:05 MST"}

	// initialize logger
	if err := logger.Init(cfg.LogLevel, cfg.LogTimeFormat); err != nil {
		return fmt.Errorf("failed to initialize logger: %v", err)
	}

	// Load .env configuration
	err := godotenv.Load()
	if err != nil {
		logger.Log.Warn(".env file not found, using environment variables")
	}

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URL"))
	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("error getting connect mongo client: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	utility.CrawlDaily()

	// initialize model
	model := model2.InitModel(mongoClient)

	// initialize handlers
	handler := handlers.NewHandlers(model)

	// initialize scheduler
	go func() {
		it := utility.Scheduler{Enabled: true, Job: model.DisableInactiveUsers}
		it.Start()
	}()
	go func() {
		it := utility.DailyScheduler{Enabled: true, Job: model.UpdateDailies}
		it.Start()
	}()

	// run HTTP gateway
	go func() {
		_ = rest.RunServer(ctx, cfg.GRPCPort, cfg.HTTPPort)
	}()

	fmt.Printf("%+v", utility.CrawlCasesByDate("2021-01-01", "2021-01-20")[0])

	return grpc.RunServer(ctx, handler, cfg.GRPCPort)
}
