package main

import (
	"database/sql"
	"log"
	"order-svc/config"
	handlers "order-svc/handlers/order"
	repository "order-svc/repository/order"
	redisSvc "order-svc/repository/redis"
	redisRepo "order-svc/repository/redis/scheduler"
	"order-svc/routes"
	usecases "order-svc/usecases/order"

	"github.com/go-redis/redis/v8"

	// Kafka broker import
	kafkaBroker "order-svc/handlers/broker/kafka"
	taskKafka "order-svc/handlers/broker/kafka/consumer/task"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		return
	}

	db, err := config.InitPostgreSQL(cfg.Postgres)
	if err != nil {
		return
	}
	defer db.Close()

	redisSv, err := redisSvc.NewRedisClient(cfg.Redis)
	if err != nil {
		return
	}

	kafkaInstance, err := kafkaBroker.NewKafkaConsumer(cfg.Kafka) // Ganti dengan alamat broker Kafka Anda
	if err != nil {
		log.Fatalf("Failed to initialize Kafka: %s", err)
	}

	routes := initDepedencies(db, redisSv, kafkaInstance)
	routes.Setup(cfg.BaseURL)
	routes.Run(cfg.Port)
}

func initDepedencies(db *sql.DB, rd *redis.Client, k *kafkaBroker.KafkaConsumer) *routes.Routes {
	orderRepo := repository.NewOrderRepository(db)
	redisRepo := redisRepo.NewBookingSchedulerService(rd, orderRepo)
	// / Start Redis Worker in a Goroutine
	go func() {
		redisRepo.StartWorker()
	}()
	orderUsecase := usecases.NewOrderUsecase(orderRepo, redisRepo)

	topics := "payOrder"
	taskWorker := taskKafka.NewTaskWorker(k, orderUsecase, topics)
	go taskWorker.InitKafka()
	orderHandler := handlers.NewOrderHandler(orderUsecase)

	return &routes.Routes{
		OrderHandler: orderHandler,
	}
}
