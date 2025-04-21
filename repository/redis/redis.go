package redis

import (
	"context"
	"log"
	"order-svc/config"
	"time"

	"github.com/go-redis/redis/v8" // Import Redis Go Client v8
)

// NewRedisClient membuat koneksi ke Redis dan mengembalikan *redis.Client
func NewRedisClient(conf config.RedisConfig) (*redis.Client, error) {
	ctx := context.Background()

	// Membuat instance Redis client dengan konfigurasi dari file config
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Address, // Alamat host dan port Redis
		Password: "",           // Password Redis (jika ada, kosong jika tidak di-set)

		// Konfigurasi Connection Pool
		PoolSize:     1,                 // Maksimum jumlah koneksi dalam pool
		MinIdleConns: 1,                 // Jumlah minimum koneksi idle
		IdleTimeout:  time.Duration(60), // Waktu maksimum koneksi idle sebelum ditutup
	})

	// Melakukan ping untuk memastikan koneksi Redis berhasil
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Cannot connect to Redis: %s", err)
		return nil, err // Jika gagal, kembalikan error agar bisa ditangani aplikasi
	}

	log.Print("Connected to Redis successfully!") // Logging jika berhasil
	return client, nil
}
