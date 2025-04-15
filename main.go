package main

import (
	"database/sql"
	"order-svc/config"
	handlers "order-svc/handlers/user"
	repository "order-svc/repository/user"
	"order-svc/routes"
	usecases "order-svc/usecases/user"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
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

	redis, err := config.InitRedis(cfg.Redis)
	if err != nil {
		return
	}
	defer redis.Close()

	rpc, err := config.RPCDial(cfg.Grpc)
	if err != nil {
		return
	}

	routes := initDepedencies(db, rpc, redis)
	routes.Setup(cfg.BaseURL)
	routes.Run(cfg.Port)
}

func initDepedencies(db *sql.DB, rpc *grpc.ClientConn, redis *redis.Client) *routes.Routes {
	userRepo := repository.NewUserStore(db)
	userUC := usecases.NewUserUsecase(userRepo, redis)
	userHandler := handlers.NewUserHandler(userUC)

	return &routes.Routes{
		User: userHandler,
	}
}
