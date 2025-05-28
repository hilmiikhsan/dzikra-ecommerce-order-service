package cmd

import (
	"flag"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/constants"
	externalNotification "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/notification"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/infrastructure/config"
	redisRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/infrastructure/redis"
	orderConsumer "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/integration/rabbitmq/consumer/ports"
	consumerSerice "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/integration/rabbitmq/consumer/service"
	deliveryConsumer "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/integration/rabbitmq/delivery_consumer/ports"
	deliveryConsumerSerice "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/integration/rabbitmq/delivery_consumer/service"
	rajaongkirService "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/integration/rajaongkir/service"
	orderRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order/repository"
	orderItemHistoryRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_item_history/repository"
	orderPaymentRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_payment/repository"
	orderStatusHistory "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_status_history/repository"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/route"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/rs/zerolog/log"
)

func RunServerHTTP(cmd *flag.FlagSet, args []string) {
	var (
		envs        = config.Envs
		flagAppPort = cmd.String("port", envs.App.Port, "Application port")
		SERVER_PORT string
	)

	// logLevel, err := zerolog.ParseLevel(envs.App.LogLevel)
	// if err != nil {
	// 	logLevel = zerolog.InfoLevel
	// }

	if err := cmd.Parse(args); err != nil {
		log.Fatal().Err(err).Msg("Error while parsing flags")
	}

	if envs.App.Port != "" {
		SERVER_PORT = envs.App.Port
	} else {
		SERVER_PORT = *flagAppPort
	}

	app := fiber.New()

	// Application Middlewares
	if envs.App.Environtment == constants.EnvProduction {
		app.Use(limiter.New(limiter.Config{
			Max:        50,
			Expiration: 30 * time.Second,
		}))
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS,HEAD",
		AllowHeaders: "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization",
	}))
	// End Application Middlewares

	adapter.Adapters.Sync(
		adapter.WithRestServer(app),
		adapter.WithDzikraPostgres(),
		adapter.WithDzikraRedis(),
		adapter.WithDzikraMidtrans(),
		adapter.WithRabbitMQ(),
		adapter.WithValidator(validator.NewValidator()),
	)

	// infrastructure.InitializeLogger(
	// 	envs.App.Environtment,
	// 	envs.App.LogFile,
	// 	logLevel,
	// )

	app.Get("/metrics", monitor.New(monitor.Config{Title: config.Envs.App.Name + config.Envs.App.Environtment + " Metrics"}))
	route.SetupRoutes(app)

	// print all routes that are registered
	// for _, route := range app.Stack() {
	// 	for _, handler := range route {
	// 		fmt.Printf("Method: %s, Path: %s\n", handler.Method, handler.Path)
	// 	}
	// }

	// repository
	orderRepository := orderRepository.NewOrderRepository(adapter.Adapters.DzikraPostgres)
	orderStatusHistoryRepository := orderStatusHistory.NewOrderStatusHistoryRepository(adapter.Adapters.DzikraPostgres)
	orderPaymentRepository := orderPaymentRepository.NewOrderPaymentRepository(adapter.Adapters.DzikraPostgres)
	externalNotification := &externalNotification.External{}
	redisRepository := redisRepository.NewRedisRepository(adapter.Adapters.DzikraRedis)
	rajaongkirService := rajaongkirService.NewRajaongkirService(redisRepository)
	orderItemHistoryRepository := orderItemHistoryRepository.NewOrderItemHistoryRepository(adapter.Adapters.DzikraPostgres)

	// consumer expire service
	consumerSerice := consumerSerice.NewConsumerService(
		adapter.Adapters.DzikraPostgres,
		orderRepository,
		orderPaymentRepository,
		orderStatusHistoryRepository,
		externalNotification,
	)

	// Start expire consumer
	go func() {
		if err := orderConsumer.StartExpireConsumer(
			adapter.Adapters.RabbitMQConn,
			adapter.Adapters.DzikraPostgres,
			adapter.Adapters.DzikraRedis,
			consumerSerice,
		); err != nil {
			log.Fatal().Err(err).Msg("Failed to start expire consumer")
		}
	}()

	// delivery consumer service
	deliveryConsumerSerice := deliveryConsumerSerice.NewDeliveryConsumerService(
		adapter.Adapters.DzikraPostgres,
		orderRepository,
		orderPaymentRepository,
		orderStatusHistoryRepository,
		externalNotification,
		rajaongkirService,
		orderItemHistoryRepository,
	)

	// Start delivery consumer
	go func() {
		if err := deliveryConsumer.StartDeliveryConsumer(
			adapter.Adapters.RabbitMQConn,
			adapter.Adapters.DzikraPostgres,
			adapter.Adapters.DzikraRedis,
			deliveryConsumerSerice,
		); err != nil {
			log.Fatal().Err(err).Msg("start delivery consumer")
		}
	}()

	// Run server in goroutine
	go func() {
		log.Info().Msgf("Server is running on port %s", SERVER_PORT)
		if err := app.Listen(":" + SERVER_PORT); err != nil {
			log.Fatal().Msgf("Error while starting server: %v", err)
		}
	}()
	// End Run server in goroutine

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)

	shutdownSignals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}
	if runtime.GOOS == "windows" {
		shutdownSignals = []os.Signal{os.Interrupt}
	}

	signal.Notify(quit, shutdownSignals...)
	<-quit
	log.Info().Msg("Server is shutting down ...")

	err := adapter.Adapters.Unsync()
	if err != nil {
		log.Error().Msgf("Error while closing adapters: %v", err)
	}

	log.Info().Msg("Server gracefully stopped")
}
