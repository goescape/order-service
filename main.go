package main

import (
	"database/sql"
	"order-svc/config"
	handlers "order-svc/handlers/user"
	repository "order-svc/repository/user"
	"order-svc/routes"
	usecases "order-svc/usecases/user"
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

	routes := initDepedencies(db)
	routes.Setup(cfg.BaseURL)
	routes.Run(cfg.Port)
}

func initDepedencies(db *sql.DB) *routes.Routes {
	userRepo := repository.NewUserStore(db)
	userUC := usecases.NewUserUsecase(userRepo)
	userHandler := handlers.NewUserHandler(userUC)

	return &routes.Routes{
		User: userHandler,
	}
}
