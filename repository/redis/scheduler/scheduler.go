package scheduler

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	dto "order-svc/model"
	repo "order-svc/repository/order"

	"github.com/go-redis/redis/v8"
)

const RedisExpiredEvent = "__keyevent@*__:expired"

// Event __keyevent@*__:expired adalah nama khusus yang digunakan oleh Redis untuk keyspace notifications.
// Ini adalah bagian dari mekanisme bawaan Redis untuk memberi tahu sistem lain saat suatu kunci (key)
// di Redis telah kedaluwarsa (expired).

// Kenapa Nama Eventnya Seperti Itu?
// Nama ini terdiri dari beberapa bagian:

// __keyevent@0__:
// keyevent → Menandakan bahwa ini adalah event terkait perubahan pada suatu key.

// @* → Menunjukkan bahwa event ini terjadi di database Redis dengan indeks *
// (karena Redis bisa memiliki beberapa database, default-nya adalah 0)

// expired:
// Menunjukkan bahwa event ini akan dipublikasikan ketika suatu key di Redis kedaluwarsa (karena TTL-nya habis).

// Redis tidak memungkinkan kita untuk mengubah nama event ini karena ini adalah bagian dari
// internal keyspace notifications yang sudah ditentukan oleh Redis sendiri.

// Interface untuk scheduler booking
type SchedulerInterface interface {
	ScheduleTaskCancellation(orderID string) error
	StartWorker() // Memulai worker untuk mendengarkan event Redis
}

// Struct implementasi scheduler
type bookingSchedulerService struct {
	redisClient *redis.Client        // Redis client untuk menyimpan TTL booking
	Repo        repo.OrderRepository // Repository untuk akses database booking
}

// Constructor untuk membuat service scheduler
func NewBookingSchedulerService(redisClient *redis.Client, r repo.OrderRepository) SchedulerInterface {
	return &bookingSchedulerService{
		redisClient: redisClient,
		Repo:        r,
	}
}

func (s *bookingSchedulerService) ScheduleTaskCancellation(orderID string) error {
	ctx := context.Background()
	key := fmt.Sprintf("task:%s:expire", orderID) // Format key unik untuk Redis

	ttl := 60 * time.Second

	// Menyimpan key di Redis dengan TTL sekian waktu
	err := s.redisClient.SetEX(ctx, key, orderID, ttl).Err()
	if err != nil {
		log.Println("Gagal menjadwalkan pembatalan task:", err)
		return err
	}

	log.Printf("Task ID %s dijadwalkan untuk dibatalkan dalam %.2f detik", orderID, ttl.Seconds())
	return nil
}

// Worker yang berjalan terus-menerus untuk mendengarkan event expiration dari Redis
// func (s *bookingSchedulerService) StartWorker() {
// 	ctx := context.Background()

// 	// Aktifkan keyspace notifications jika belum aktif
// 	err := s.redisClient.ConfigSet(ctx, "notify-keyspace-events", "Ex").Err()
// 	if err != nil {
// 		log.Println("Gagal mengaktifkan keyspace notifications:", err)
// 		return
// 	}

// 	pubsub := s.redisClient.PSubscribe(ctx, RedisExpiredEvent) // Subscribe ke event Redis expiration

// 	log.Println("Worker Redis berjalan... Mendengarkan event expired")

// 	for {
// 		// Menerima pesan dari Redis ketika ada key yang expired
// 		msg, err := pubsub.ReceiveMessage(ctx)
// 		if err != nil {
// 			log.Println("Error menerima pesan Redis:", err)
// 			continue
// 		}

// 		// Mengambil Task ID dari key yang expired
// 		var data dto.CancelOrderModel
// 		_, err = fmt.Sscanf(msg.Payload, "task:%d:expire", &data.ID)
// 		if err != nil {
// 			continue
// 		}

// 		// Membatalkan booking karena tidak dibayar dalam waktu yang ditentukan
// 		log.Printf("Membatalkan task ID %s karena tidak diselesaikan", data.ID)
// 		err = s.Repo.CancelOrder(ctx, &data)
// 		if err != nil {
// 			log.Println("Gagal membatalkan task:", err)
// 		} else {
// 			log.Printf("Task ID %s berhasil dibatalkan", data.ID)
// 		}
// 	}
// }

func (s *bookingSchedulerService) StartWorker() {
	ctx := context.Background()
	pubsub := s.redisClient.PSubscribe(ctx, "__keyevent@0__:expired") // Gunakan DB 0 secara eksplisit

	log.Println("Worker Redis aktif... Menunggu key expired")

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			log.Println("Error menerima pesan Redis:", err)
			continue
		}

		log.Println("Pesan Redis diterima:", msg.Channel, msg.Payload) // ✅ Debug pesan

		// Contoh payload: "task:12345:expire"
		var data dto.CancelOrderModel
		payload := msg.Payload

		// Cek dan potong prefix + suffix
		if strings.HasPrefix(payload, "task:") && strings.HasSuffix(payload, ":expire") {
			orderID := strings.TrimSuffix(strings.TrimPrefix(payload, "task:"), ":expire")
			data.ID = orderID

			log.Printf("Membatalkan task ID %s karena tidak diselesaikan", data.ID)
			err = s.Repo.CancelOrder(ctx, &data)
			if err != nil {
				log.Println("Gagal membatalkan task:", err)
			} else {
				log.Printf("Task ID %s berhasil dibatalkan", data.ID)
			}
		} else {
			log.Println("Format payload tidak sesuai:", payload)
		}
	}
}
