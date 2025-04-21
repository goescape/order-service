package kafka

import (
	"log"
	"order-svc/config"

	"github.com/IBM/sarama"
)

// KafkaConsumer adalah struct yang menyimpan instance Kafka consumer
type KafkaConsumer struct {
	Consumer sarama.Consumer
}

// NewKafkaConsumer menginisiasi Kafka consumer berdasarkan config dan mengembalikan koneksi Kafka
func NewKafkaConsumer(cfg config.KafkaConf) (*KafkaConsumer, error) {
	// Set konfigurasi untuk Kafka
	saramaConfig := sarama.NewConfig()

	// Buat Kafka consumer
	consumer, err := sarama.NewConsumer([]string{cfg.Broker}, saramaConfig)
	if err != nil {

		return nil, err
	}

	log.Printf("Kafka consumer berhasil terkoneksi ke: %s", cfg.Broker)

	// Mengembalikan pointer ke KafkaConsumer
	return &KafkaConsumer{Consumer: consumer}, nil
}
