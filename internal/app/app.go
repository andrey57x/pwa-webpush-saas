package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/andrey57x/pwa-webpush-saas/internal/config"
	"github.com/andrey57x/pwa-webpush-saas/internal/connectors/kafka"
	"github.com/andrey57x/pwa-webpush-saas/internal/connectors/postgres"
	"github.com/andrey57x/pwa-webpush-saas/internal/connectors/redis"
	"github.com/andrey57x/pwa-webpush-saas/internal/delivery"
	"github.com/andrey57x/pwa-webpush-saas/internal/pkg/pushsender"
	repoPg "github.com/andrey57x/pwa-webpush-saas/internal/repository/postgres"
	repoRedis "github.com/andrey57x/pwa-webpush-saas/internal/repository/redis"
	"github.com/andrey57x/pwa-webpush-saas/internal/router"
	"github.com/andrey57x/pwa-webpush-saas/internal/usecase"
	"github.com/andrey57x/pwa-webpush-saas/migrations"
	redisDriver "github.com/redis/go-redis/v9"
)

type App struct {
	cfg         *config.Config
	sqlDB       *sql.DB
	redisClient *redisDriver.Client
	producer    *kafka.Producer
	consumer    *kafka.Consumer
	server      *http.Server
}

func New(cfg *config.Config) (*App, error) {
	ctx := context.Background()

	if err := migrations.RunMigrations(cfg.GetPostgresDSN()); err != nil {
		return nil, fmt.Errorf("migration error: %w", err)
	}

	sqlDB, err := postgres.NewPool(cfg.GetPostgresDSN())
	if err != nil {
		return nil, fmt.Errorf("postgres init error: %w", err)
	}

	redisClient, err := redis.NewClient(ctx, cfg.GetRedisAddr(), cfg.RedisPassword)
	if err != nil {
		return nil, fmt.Errorf("redis init error: %w", err)
	}

	subRepo := repoPg.NewSubscriptionRepository(sqlDB)
	logRepo := repoPg.NewDeliveryLogRepository(sqlDB)
	appRepo := repoPg.NewAppRepository(sqlDB)
	campaignRepo := repoPg.NewCampaignRepository(sqlDB)
	dedupRepo := repoRedis.NewDedupRepository(redisClient)

	brokers := strings.Split(cfg.KafkaBrokers, ",")
	producer := kafka.NewProducer(brokers, cfg.KafkaTopicPushJobs)

	sender := pushsender.NewSender(cfg.ContactEmail)
	workerUsecase := usecase.NewWorkerUsecase(subRepo, logRepo, sender)

	consumer := kafka.NewConsumer(
		brokers,
		cfg.KafkaTopicPushJobs,
		cfg.KafkaConsumerGroup,
		cfg.WorkerPoolSize,
		workerUsecase,
	)

	subUsecase := usecase.NewSubscriptionUsecase(subRepo, appRepo)
	appUsecase := usecase.NewAppUsecase(appRepo)
	campaignUsecase := usecase.NewCampaignUsecase(campaignRepo, subRepo, appRepo, dedupRepo, producer)
	feedbackUsecase := usecase.NewFeedbackUsecase(logRepo)

	subHandler := delivery.NewSubscriptionHandler(subUsecase)
	appHandler := delivery.NewAppHandler(appUsecase)
	campaignHandler := delivery.NewCampaignHandler(campaignUsecase, campaignRepo)
	feedbackHandler := delivery.NewFeedbackHandler(feedbackUsecase)

	r := router.NewRouter(subHandler, campaignHandler, appHandler, feedbackHandler)

	httpServer := &http.Server{
		Addr:    ":" + cfg.APIPort,
		Handler: r,
	}

	return &App{
		cfg:         cfg,
		sqlDB:       sqlDB,
		redisClient: redisClient,
		producer:    producer,
		consumer:    consumer,
		server:      httpServer,
	}, nil
}

func (a *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	go a.consumer.Start(ctx)

	go func() {
		log.Printf("HTTP Server listening on http://localhost:%s", a.cfg.APIPort)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP Server failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = a.server.Shutdown(shutdownCtx)
	_ = a.consumer.Close()
	_ = a.producer.Close()
	_ = a.redisClient.Close()
	_ = a.sqlDB.Close()

	log.Println("Application stopped cleanly.")
	return nil
}
