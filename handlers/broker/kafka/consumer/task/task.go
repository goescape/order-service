package task

import (
	"context"
	"encoding/json"
	"log"

	kafkaBroker "order-svc/handlers/broker/kafka" // Kafka Broker
	"order-svc/model"

	useCase "order-svc/usecases/order"

	"github.com/IBM/sarama"
)

// Interface untuk inisialisasi Kafka Consumer
type NotifTaskInterface interface {
	InitKafka()
}

// Struct untuk worker yang menangani task dari Kafka
type TaskWorkerImpl struct {
	kafka   *kafkaBroker.KafkaConsumer // Instance Kafka connection
	topics  string                     // Nama topik Kafka
	UseCase useCase.OrderUsecases      // Use case untuk task
}

// Konstruktor untuk membuat TaskWorker
func NewTaskWorker(kafka *kafkaBroker.KafkaConsumer, useCase useCase.OrderUsecases, topics string) NotifTaskInterface {
	taskWorkerImpl := &TaskWorkerImpl{
		kafka:   kafka,
		topics:  topics, // Anda bisa mendefinisikan topik Kafka sesuai kebutuhan
		UseCase: useCase,
	}

	return taskWorkerImpl
}

// Fungsi untuk inisialisasi Kafka consumer
// Fungsi untuk inisialisasi Kafka consumer
func (p *TaskWorkerImpl) InitKafka() {
	// Buat Kafka consumer
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil) // Ganti dengan broker Kafka Anda
	if err != nil {
		log.Fatal("Error creating Kafka consumer: ", err)
	}

	// Proses setiap topik yang ada
	go p.consumeTopic(consumer, p.topics)

}

// Fungsi untuk menangani konsumsi pesan dari topik Kafka
func (p *TaskWorkerImpl) consumeTopic(consumer sarama.Consumer, topic string) {
	ctx := context.Background()
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error starting partition consumer for topic %s: %v", topic, err)
	}
	defer partitionConsumer.Close()

	// Proses pesan yang diterima dari topik
	for msg := range partitionConsumer.Messages() {
		log.Printf("Received message on topic %s: %s", topic, string(msg.Value))

		// Menangani pesan sesuai dengan topiknya

		// Menangani Pay Task
		taskDTO := &model.PayOrderModel{}
		if err := json.Unmarshal(msg.Value, &taskDTO); err != nil {
			log.Printf("Error parsing Pay Task payload: %+v", err)
			continue
		}
		if err := p.UseCase.PayOrder(ctx, taskDTO); err != nil {
			log.Printf("Error executing Pay Task: %+v", err)
		}

	}
}
