package repository

import (
	"context"
	"database/sql"
	"order-svc/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, req *model.CreateOrderReq) (*model.CreateOrderResp, error)
	GetOrderList(ctx context.Context, req *model.GetOrderListRequest) (*model.ListOrderResponse, error)
	CancelOrder(ctx context.Context, req *model.CancelOrderModel) error
	PayOrder(ctx context.Context, req *model.PayOrderModel) error
}

type orderStore struct {
	db *sql.DB
}

var _ OrderRepository = &orderStore{}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderStore{
		db: db,
	}
}
