package main

import (
	"database/sql"
	"order-svc/config"
	handlers "order-svc/handlers/order"
	repository "order-svc/repository/order"
	"order-svc/routes"
	usecases "order-svc/usecases/order"
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
	orderRepo := repository.NewOrderRepository(db)
	orderUsecase := usecases.NewOrderUsecase(orderRepo)
	orderHandler := handlers.NewOrderHandler(orderUsecase)

	return &routes.Routes{
		OrderHandler: orderHandler,
	}
}
